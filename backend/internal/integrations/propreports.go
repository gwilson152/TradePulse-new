package integrations

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tradepulse/api/internal/models"
)

type PropReportsClient struct {
	BaseURL  string
	Username string
	Password string
	client   *http.Client
	token    string
}

// PropReports CSV fill record format
type PropReportsFill struct {
	DateTime   string // Date/Time
	Account    string // Account
	Side       string // B/S
	Qty        string // Qty
	Symbol     string // Symbol
	Price      string // Price
	Route      string // Route
	Liq        string // Liq
	Comm       string // Comm
	EcnFee     string // Ecn Fee
	SEC        string // SEC
	TAF        string // TAF
	NSCC       string // NSCC
	Clr        string // Clr
	Misc       string // Misc
	OrderId    string // Order Id
	FillId     string // Fill Id
	Currency   string // Currency
	ISIN       string // ISIN
	CUSIP      string // CUSIP
	Status     string // Status
	PropReportsId string // PropReports Id
}

func NewPropReportsClient(site, username, password string) *PropReportsClient {
	baseURL := fmt.Sprintf("https://%s", site)

	return &PropReportsClient{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Login authenticates with PropReports and stores the token
func (c *PropReportsClient) Login() error {
	apiURL := fmt.Sprintf("%s/api.php", c.BaseURL)

	data := url.Values{}
	data.Set("action", "login")
	data.Set("user", c.Username)
	data.Set("password", c.Password)

	resp, err := c.client.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	token := strings.TrimSpace(string(body))
	if token == "" {
		return fmt.Errorf("received empty token from PropReports")
	}

	c.token = token
	return nil
}

// Logout expires the token
func (c *PropReportsClient) Logout() error {
	if c.token == "" {
		return nil
	}

	apiURL := fmt.Sprintf("%s/api.php", c.BaseURL)
	data := url.Values{}
	data.Set("action", "logout")
	data.Set("token", c.token)

	resp, err := c.client.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	defer resp.Body.Close()

	c.token = ""
	return nil
}

// FetchTrades fetches trades (fills) from PropReports
func (c *PropReportsClient) FetchTrades(fromDate, toDate string) ([]models.Trade, error) {
	// Login first
	if err := c.Login(); err != nil {
		return nil, err
	}
	defer c.Logout()

	// Get account ID (we'll query all accounts with groupId=-2)
	apiURL := fmt.Sprintf("%s/api.php", c.BaseURL)

	// Set date range defaults if not provided
	if fromDate == "" {
		fromDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02") // Last month
	}
	if toDate == "" {
		toDate = time.Now().Format("2006-01-02") // Today
	}

	data := url.Values{}
	data.Set("action", "fills")
	data.Set("token", c.token)
	data.Set("groupId", "-2") // All accounts
	data.Set("startDate", fromDate)
	data.Set("endDate", toDate)
	data.Set("page", "1")

	resp, err := c.client.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fills: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fills request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse CSV response
	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV response: %w", err)
	}

	if len(records) < 2 {
		return []models.Trade{}, nil // No data (only header or empty)
	}

	// Convert fills to trades
	trades, err := c.convertFillsToTrades(records[1:]) // Skip header
	if err != nil {
		return nil, err
	}

	return trades, nil
}

// convertFillsToTrades converts PropReports fill records (CSV rows) into Trade models
// PropReports returns individual fills, but we need to group them into trades (positions)
func (c *PropReportsClient) convertFillsToTrades(records [][]string) ([]models.Trade, error) {
	// Group fills by symbol + order to reconstruct trades
	type TradeKey struct {
		Symbol string
		Date   string
	}

	fillsBySymbol := make(map[TradeKey][]PropReportsFill)

	for _, record := range records {
		if len(record) < 7 {
			continue // Skip incomplete records
		}

		// Skip "Page X/Y" trailer lines
		if strings.HasPrefix(record[0], "Page ") {
			continue
		}

		fill := PropReportsFill{
			DateTime: record[0],
			Account:  record[1],
			Side:     record[2],
			Qty:      record[3],
			Symbol:   record[4],
			Price:    record[5],
		}

		// Extract additional fields if present
		if len(record) > 6 {
			fill.Route = record[6]
		}
		if len(record) > 8 {
			fill.Comm = record[8]
		}

		// Parse date to group by day
		dateTime, err := time.Parse("01/02/2006 15:04:05", fill.DateTime)
		if err != nil {
			continue // Skip fills with invalid timestamps
		}

		key := TradeKey{
			Symbol: fill.Symbol,
			Date:   dateTime.Format("2006-01-02"),
		}

		fillsBySymbol[key] = append(fillsBySymbol[key], fill)
	}

	// Convert grouped fills into trades
	trades := make([]models.Trade, 0)

	for key, fills := range fillsBySymbol {
		if len(fills) == 0 {
			continue
		}

		// Determine if this is a complete round-trip trade
		var buys, sells []PropReportsFill
		for _, fill := range fills {
			if fill.Side == "B" {
				buys = append(buys, fill)
			} else if fill.Side == "S" || fill.Side == "T" {
				sells = append(sells, fill)
			}
		}

		// Calculate totals
		var buyQty, sellQty float64
		var buyTotal, sellTotal float64
		var totalFees float64

		for _, buy := range buys {
			qty, _ := strconv.ParseFloat(buy.Qty, 64)
			price, _ := strconv.ParseFloat(buy.Price, 64)
			comm, _ := strconv.ParseFloat(buy.Comm, 64)

			buyQty += qty
			buyTotal += qty * price
			totalFees += comm
		}

		for _, sell := range sells {
			qty, _ := strconv.ParseFloat(sell.Qty, 64)
			price, _ := strconv.ParseFloat(sell.Price, 64)
			comm, _ := strconv.ParseFloat(sell.Comm, 64)

			sellQty += qty
			sellTotal += qty * price
			totalFees += comm
		}

		// Determine trade type and entry
		var tradeType models.TradeType
		var entryPrice, exitPrice float64
		var quantity float64
		var openTime, closeTime time.Time

		if len(buys) > 0 && len(sells) > 0 {
			// Complete round trip
			if buys[0].DateTime < sells[0].DateTime {
				// Long trade (buy first, sell later)
				tradeType = models.TradeLong
				entryPrice = buyTotal / buyQty
				exitPrice = sellTotal / sellQty
				quantity = buyQty
				openTime, _ = time.Parse("01/02/2006 15:04:05", buys[0].DateTime)
				closeTime, _ = time.Parse("01/02/2006 15:04:05", sells[len(sells)-1].DateTime)
			} else {
				// Short trade (sell first, buy later)
				tradeType = models.TradeShort
				entryPrice = sellTotal / sellQty
				exitPrice = buyTotal / buyQty
				quantity = sellQty
				openTime, _ = time.Parse("01/02/2006 15:04:05", sells[0].DateTime)
				closeTime, _ = time.Parse("01/02/2006 15:04:05", buys[len(buys)-1].DateTime)
			}

			// Calculate P&L
			var pnl float64
			if tradeType == models.TradeLong {
				pnl = (exitPrice - entryPrice) * quantity - totalFees
			} else {
				pnl = (entryPrice - exitPrice) * quantity - totalFees
			}

			trade := models.Trade{
				Symbol:     key.Symbol,
				TradeType:  tradeType,
				Quantity:   quantity,
				EntryPrice: entryPrice,
				ExitPrice:  &exitPrice,
				OpenedAt:   openTime,
				ClosedAt:   &closeTime,
				PnL:        &pnl,
				Fees:       totalFees,
			}

			trades = append(trades, trade)
		} else if len(buys) > 0 {
			// Open long position
			tradeType = models.TradeLong
			entryPrice = buyTotal / buyQty
			quantity = buyQty
			openTime, _ = time.Parse("01/02/2006 15:04:05", buys[0].DateTime)

			trade := models.Trade{
				Symbol:     key.Symbol,
				TradeType:  tradeType,
				Quantity:   quantity,
				EntryPrice: entryPrice,
				OpenedAt:   openTime,
				Fees:       totalFees,
			}

			trades = append(trades, trade)
		} else if len(sells) > 0 {
			// Open short position
			tradeType = models.TradeShort
			entryPrice = sellTotal / sellQty
			quantity = sellQty
			openTime, _ = time.Parse("01/02/2006 15:04:05", sells[0].DateTime)

			trade := models.Trade{
				Symbol:     key.Symbol,
				TradeType:  tradeType,
				Quantity:   quantity,
				EntryPrice: entryPrice,
				OpenedAt:   openTime,
				Fees:       totalFees,
			}

			trades = append(trades, trade)
		}
	}

	return trades, nil
}

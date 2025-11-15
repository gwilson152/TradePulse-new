# CSV Import Guide

TradePulse supports importing your trading history from CSV files exported by various trading platforms. This guide explains how to prepare and import your trade data.

## Supported Platforms

### DAS Trader Pro

DAS Trader Pro is currently the primary supported platform. The import system can handle DAS Trader's execution-level exports and automatically reconstruct complete positions.

**Export Instructions:**

1. Open DAS Trader Pro
2. Click **Trade** menu → **Trade Log**
3. Right-click anywhere in the Trade Log window
4. Select **Export** to download your trades as CSV
5. **Important**: Export your trades before 10 PM ET each day, as DAS Trader only exports the current day's trades

**Important Notes:**
- DAS Trader CSV exports only include **time** (HH:MM:SS), not the date
- You will need to manually enter the trading date when importing
- The import system groups individual executions into complete positions
- Make sure to export the correct trading day's data

**CSV Format:**
```
Time,Symb,Side,Price,Qty,Route,Broker,Account,Type,Cloid
09:30:15,AAPL,B,150.25,100,ARCA,,,Limit,
09:31:45,AAPL,S,150.50,100,ARCA,,,Limit,
```

**Column Mapping:**
- `Symb` or `Symbol` → Stock symbol
- `Side` → B (Buy) or S (Sell)
- `Price` or `Exec Price` → Execution price
- `Qty` or `Quantity` → Share quantity
- `Time` → Execution time (requires date input during import)
- `Commission` or `Comm` → Trading fees (optional)
- `Account` → Account identifier (optional)
- `Type` → Order type (optional)

### PropReports

PropReports exports are also supported and include more complete data.

**Export Instructions:**

1. Log in to PropReports
2. Navigate to your trade history or detailed reports
3. Export trades as CSV
4. PropReports exports include full date/time, so no manual date entry is required

**CSV Format:**

PropReports has a flexible format. The import system supports multiple possible column names for each field:

- **Symbol**: Symbol, Ticker, Instrument
- **Side**: Side, Direction, Type (values: LONG, SHORT, BUY, SELL, B, S)
- **Quantity**: Quantity, Qty, Shares, Size
- **Price**: Price, Entry Price, Avg Price, Average Price
- **Timestamp**: Date, Time, DateTime, Timestamp, Entry Time
- **Fees**: Commission, Comm, Fees, Total Fees
- **Account**: Account, Account Number

## Import Process

### Step 1: Upload

1. Navigate to the **Trades** page
2. Click the **Import CSV** button
3. Select your trading platform from the dropdown
   - Choose "DAS Trader Pro" for DAS exports
   - Choose "PropReports" for PropReports exports
4. If using DAS Trader, enter the trading date for your export
5. Drag and drop your CSV file or click to browse
6. Click **Parse CSV** to validate the file

### Step 2: Preview

After parsing, you'll see:

- **Import Statistics**
  - Total Rows: Number of rows in the CSV
  - Valid Trades: Number of successfully parsed trades
  - Warnings: Non-critical issues found
  - Errors: Critical issues that will cause rows to be skipped

- **Errors Section** (if any)
  - Lists specific issues with row numbers
  - Common errors: missing columns, invalid values, malformed data

- **Warnings Section** (if any)
  - Possible duplicate trades
  - Other non-critical issues

- **Trade Preview**
  - Shows first 10 trades that will be imported
  - Displays symbol, type, entry/exit prices, quantity, and P&L

Review the preview carefully before proceeding.

### Step 3: Import

1. Review the statistics and trade preview
2. If errors exist, trades with errors will be skipped
3. Click **Import Trades** to complete the import
4. Wait for confirmation message
5. Your trades will appear in the trades table

## How It Works

### For DAS Trader (Execution Grouping)

DAS Trader exports individual executions, not complete positions. The import system:

1. **Groups by Symbol**: All executions for the same symbol are grouped together
2. **Sorts by Time**: Executions are chronologically ordered
3. **Detects Positions**: Tracks buy/sell quantities to determine when positions open and close
4. **Determines Direction**:
   - First execution is B (buy) → LONG position
   - First execution is S (sell) → SHORT position
5. **Calculates Averages**:
   - Average entry price = weighted average of entry executions
   - Average exit price = weighted average of exit executions
6. **Computes P&L**:
   - LONG: (avg exit price - avg entry price) × shares - fees
   - SHORT: (avg entry price - avg exit price) × shares - fees
7. **Handles Multiple Round Trips**: If you traded the same symbol multiple times in one day, each complete position is treated as a separate trade

**Example:**
```
Time      Symbol  Side  Price   Qty
09:30:15  AAPL    B     150.00  50   → Open LONG position
09:30:45  AAPL    B     150.25  50   → Add to position
09:31:30  AAPL    S     150.75  100  → Close position (avg entry: $150.125)

Result: 1 LONG trade
- Entry: $150.125 (average of two buys)
- Exit: $150.75
- Quantity: 100 shares
- P&L: ($150.75 - $150.125) × 100 = +$62.50
```

### For PropReports (Position-Level Data)

PropReports already exports complete positions, so no grouping is needed. Each row becomes one trade.

## Validation and Error Handling

### Automatic Validation

The import system validates:

- **File Type**: Must be .csv
- **File Size**: Maximum 10 MB
- **Required Columns**: Symbol, Side, Price, Quantity, Timestamp must be present
- **Data Types**: Prices and quantities must be valid numbers
- **Date/Time**: Timestamps must be parseable
- **Side Values**: Must be recognizable (B, S, BUY, SELL, LONG, SHORT)

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| Missing required columns | CSV doesn't have expected headers | Verify export format matches platform specifications |
| Invalid price | Non-numeric or negative price | Check CSV for corrupted data |
| Invalid quantity | Non-numeric or zero/negative quantity | Check CSV for corrupted data |
| Invalid timestamp | Unparseable date/time | For DAS Trader, verify time format is HH:MM:SS |
| Invalid side value | Unrecognized buy/sell indicator | Check that Side column contains B/S or BUY/SELL |
| Date is required | DAS Trader import without date | Enter the trading date in the date picker |

### Duplicate Detection

The system detects possible duplicates by comparing:
- Symbol
- Entry timestamp
- Entry price

If duplicates are found, they are flagged as warnings but will still be imported. You can manually review and delete duplicates after import.

## Best Practices

### Before Importing

1. **Export Fresh Data**: Always export your latest trading data
2. **One Day at a Time**: For DAS Trader, export and import each trading day separately
3. **Check the File**: Open the CSV in a text editor to verify it looks correct
4. **Note the Date**: Remember which trading day you're importing (critical for DAS Trader)

### During Import

1. **Review Preview**: Always check the preview before importing
2. **Check Statistics**: Ensure "Valid Trades" matches your expectations
3. **Read Errors**: If errors exist, try to fix them in the CSV and re-import
4. **Verify P&L**: Spot-check that calculated P&L values look correct

### After Import

1. **Review Trades**: Browse your imported trades in the trades table
2. **Check for Duplicates**: Look for duplicate warning flags
3. **Verify Metrics**: Check that your dashboard metrics align with your expectations
4. **Delete Duplicates**: If you accidentally imported the same data twice, delete duplicate trades

## Troubleshooting

### "Missing required columns" Error

**Problem**: The import system can't find the required CSV columns.

**Solutions**:
- Verify you selected the correct platform in the dropdown
- Open the CSV and check the header row matches the expected format
- Ensure the CSV wasn't corrupted during export or transfer
- For DAS Trader, make sure you exported from the Trade Log (not another report)

### "Invalid timestamp" Error

**Problem**: The system can't parse the time values.

**Solutions**:
- For DAS Trader: Ensure the Time column is in HH:MM:SS format
- For PropReports: Check that the date/time column is properly formatted
- Verify the trading date you entered matches the CSV data (DAS Trader only)

### Incorrect P&L Calculations

**Problem**: The P&L values don't match your expectations.

**Possible Causes**:
- **Fees not included**: Check if your CSV export includes commission/fees columns
- **Partial fills**: Verify the import captured all executions for a position
- **Wrong date**: For DAS Trader, entering the wrong date can cause timestamp issues
- **Open positions**: Positions without exits will show null P&L

**Solutions**:
- Check the entries/exits arrays in the trade detail view
- Verify all executions were included in the CSV export
- For DAS Trader, re-import with the correct date
- Open positions will need to be manually closed or updated

### Trades Not Grouping Correctly (DAS Trader)

**Problem**: Multiple executions aren't being combined into one trade.

**Possible Causes**:
- Time gaps causing the system to treat them as separate positions
- Symbol spelling variations (e.g., "AAPL" vs "AAPL ")
- Uneven buy/sell quantities creating multiple positions

**Solutions**:
- Verify the symbol is consistent throughout the CSV
- Check that buy and sell quantities match for complete positions
- Review the trade preview to see how executions were grouped
- Manual adjustment may be needed for complex scaling scenarios

### File Too Large

**Problem**: CSV file exceeds 10 MB limit.

**Solutions**:
- Export smaller date ranges
- Split large exports into multiple files
- Remove unnecessary columns (though not recommended)

## Data Privacy

- CSV files are processed **client-side** in your browser
- File contents are never sent to external servers during parsing
- Only the final parsed trade data is sent to the TradePulse API
- Your original CSV files are not stored on our servers

## Support

If you encounter issues not covered in this guide:

1. Check the error messages carefully
2. Verify your CSV format matches the platform specifications
3. Try a smaller sample file first (e.g., one or two trades)
4. Review the browser console for technical error details
5. Contact support with the specific error message and platform used

## Adding Support for New Platforms

TradePulse is designed to support additional trading platforms. If you use a platform that's not currently supported, the system can be extended by:

1. Creating a new platform schema with column mappings
2. Defining transformation functions for data parsing
3. Configuring whether execution grouping is needed

Contact the development team if you'd like to request support for a new platform.

---

**Last Updated**: November 2024
**Supported Platforms**: DAS Trader Pro, PropReports

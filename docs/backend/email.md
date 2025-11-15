# Email System

## Overview

TradePulse supports multiple email providers through an abstracted interface.

## Providers

- **SMTP** - Direct SMTP (supports trusted relay with blank credentials)
- **M365 Graph** - Microsoft 365 / Azure AD
- **Gmail** - Google Workspace / Gmail API

## Configuration

Set `EMAIL_PROVIDER` in `.env`:

```bash
EMAIL_PROVIDER=smtp  # or graph, gmail
```

### SMTP Setup

```bash
EMAIL_PROVIDER=smtp
SMTP_HOST=mail.example.com
SMTP_PORT=587
SMTP_USER=              # Can be blank for trusted relay
SMTP_PASS=              # Can be blank for trusted relay
SMTP_FROM=support@drivenw.com
```

**Trusted Relay**: Leave `SMTP_USER` and `SMTP_PASS` blank if your SMTP server allows relay from the application server's IP.

### M365 Graph Setup

```bash
EMAIL_PROVIDER=graph
GRAPH_TENANT_ID=your-tenant-id
GRAPH_CLIENT_ID=your-client-id
GRAPH_CLIENT_SECRET=your-client-secret
GRAPH_FROM_EMAIL=noreply@yourcompany.com
```

### Gmail Setup

```bash
EMAIL_PROVIDER=gmail
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REFRESH_TOKEN=your-refresh-token
GOOGLE_FROM_EMAIL=noreply@yourcompany.com
```

## Usage (To Be Implemented)

### Send Magic Link

```go
// In auth handler
func SendMagicLink(email string, token string, emailService EmailProvider) error {
    link := fmt.Sprintf("%s/auth/verify?token=%s", baseURL, token)

    subject := "Your TradePulse Login Link"
    body := fmt.Sprintf(`
        Click the link below to log in to TradePulse:

        %s

        This link expires in 15 minutes.
    `, link)

    return emailService.Send(email, subject, body)
}
```

### HTML Email

```go
htmlBody := `
<!DOCTYPE html>
<html>
<body>
    <h1>Welcome to TradePulse</h1>
    <p>Click below to log in:</p>
    <a href="{{.Link}}">Login to TradePulse</a>
    <p>This link expires in 15 minutes.</p>
</body>
</html>
`

emailService.SendHTML(email, subject, htmlBody)
```

## Email Provider Interface

```go
type EmailProvider interface {
    Send(to, subject, body string) error
    SendHTML(to, subject, htmlBody string) error
}
```

## Testing

### Test SMTP Connection

```bash
# Using swaks (SMTP test tool)
swaks --to test@example.com \
      --from support@drivenw.com \
      --server mail.example.com:587 \
      --auth LOGIN \
      --auth-user your-user \
      --auth-password your-pass
```

### Test Email Sending

```go
// In handler or test
emailService := email.NewSMTPProvider(config)
err := emailService.Send(
    "test@example.com",
    "Test Email",
    "This is a test message from TradePulse",
)
if err != nil {
    log.Printf("Failed to send email: %v", err)
}
```

## Best Practices

1. **Handle failures gracefully** - Don't block requests if email fails
2. **Log email events** - Track sent emails for debugging
3. **Use HTML templates** - For better formatting
4. **Include unsubscribe links** - For transactional emails
5. **Rate limit** - Don't spam users
6. **Test in development** - Use mailtrap.io or similar

## Troubleshooting

### SMTP Connection Failed

- Verify SMTP host and port
- Check firewall rules
- Test with telnet: `telnet mail.example.com 587`

### Authentication Failed

- Verify username/password
- Check if app-specific password is required
- Ensure trusted relay is configured if using blank credentials

### Emails Not Delivered

- Check spam folders
- Verify SPF/DKIM records
- Check email provider logs

## Implementation Status

Email providers need to be implemented in `internal/email/`:
- `provider.go` - Interface definition
- `smtp.go` - SMTP implementation
- `graph.go` - M365 Graph implementation
- `gmail.go` - Gmail API implementation

See handlers example for integration pattern.

# Backend Getting Started

## Prerequisites

- Go 1.21 or higher
- Access to PostgreSQL database at postgres1.drivenw.local:5432
- Git

## Setup

### 1. Install Dependencies

```bash
cd backend
go mod download
```

### 2. Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` and set required values:

```bash
# Generate a secure JWT secret
JWT_SECRET=$(openssl rand -base64 32)

# Update if needed
DB_HOST=postgres1.drivenw.local
DB_PORT=5432
DB_NAME=tradepulse
DB_USER=tradepulse
DB_PASSWORD=drv13llc!

# Email settings
EMAIL_PROVIDER=smtp
SMTP_HOST=your-smtp-host
SMTP_FROM=support@drivenw.com
```

### 3. Run Database Migrations

```bash
# Manual migration
psql -h postgres1.drivenw.local -p 5432 -U tradepulse -d tradepulse \
  -f migrations/001_initial_schema.up.sql

# Or migrations run automatically on server start
```

### 4. Start the Server

```bash
go run cmd/api/main.go
```

The API will start on port 9000 and be accessible at:
`https://api.tradepulse.drivenw.com`

## Verify Installation

### Check Health Endpoint

```bash
curl https://api.tradepulse.drivenw.com/health
```

Expected response:
```json
{"status":"ok"}
```

### Check Database Connection

Look for this log message on startup:
```
{"level":"INFO","msg":"Database connection established"}
{"level":"INFO","msg":"Database migrations completed"}
```

### Check Notification Bus

Look for this log message:
```
{"level":"INFO","msg":"Notification bus started"}
```

### Test Magic Link Authentication

1. Visit https://tradepulse.drivenw.com/auth/login
2. Enter your email address
3. Check the backend console logs for the magic link:
```
{"level":"INFO","msg":"Magic link generated","email":"you@example.com","token":"abc123...","link":"https://tradepulse.drivenw.com/auth/verify?token=abc123..."}
```
4. Copy the full link and visit it in your browser
5. You should be redirected to the dashboard

**Note**: Email sending is not yet implemented. Magic links are logged to the console for testing.

## Development Commands

### Run Server
```bash
go run cmd/api/main.go
```

### Build Binary
```bash
go build -o bin/api ./cmd/api
```

### Run Tests
```bash
go test ./...
```

### Format Code
```bash
go fmt ./...
```

### Check for Issues
```bash
go vet ./...
```

## Project Structure

```
backend/
├── cmd/api/              # Application entry point
│   └── main.go
├── internal/             # Private application code
│   ├── auth/            # Authentication logic
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data structures
│   ├── database/        # Database layer
│   ├── email/           # Email providers
│   ├── middleware/      # HTTP middleware
│   └── notifications/   # WebSocket system
├── migrations/          # SQL migrations
├── .env                 # Environment variables (not in git)
├── .env.example         # Environment template
├── go.mod              # Go dependencies
└── go.sum              # Dependency checksums
```

## Common Issues

### Database Connection Failed

**Error**: `Failed to connect to database`

**Solution**:
- Verify postgres1.drivenw.local is accessible
- Check credentials in `.env`
- Ensure database `tradepulse` exists

### Port Already in Use

**Error**: `bind: address already in use`

**Solution**:
```bash
# Find process using port 9000
lsof -i :9000
# or on Windows
netstat -ano | findstr :9000

# Kill the process
kill -9 <PID>
```

### Missing JWT Secret

**Error**: `JWT_SECRET environment variable is required`

**Solution**:
Generate and add to `.env`:
```bash
openssl rand -base64 32
```

## Next Steps

- Read [Structure](structure.md) to understand code organization
- See [API Patterns](api-patterns.md) for handler examples
- Check [Notifications](notifications.md) to publish notifications
- Review [Database](database.md) for schema details

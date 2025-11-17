# TradePulse

A comprehensive trading metrics and advanced journaling platform for serious traders.

## Overview

TradePulse helps traders track performance, document decision-making processes, and improve through detailed journaling with emotional state tracking and multimedia support.

### Key Features

- **Advanced Journaling**: Rich text entries with emotional state tracking, screenshot uploads, and voice notes
- **Comprehensive Metrics**: Track P&L, win rate, profit factor, and performance across all trades
- **CSV Import**: Import trading history from DAS Trader Pro, PropReports, and other platforms ([Import Guide](docs/csv-import.md))
- **Dual Authentication**: Magic Link (passwordless) OR Email/Password ([Auth Guide](docs/authentication.md))
- **Privacy-First**: Self-hosted with complete data ownership

## Technology Stack

- **Frontend**: SvelteKit 5 with TypeScript
- **Backend**: Go 1.21+ with Chi router
- **Database**: PostgreSQL 15+
- **Authentication**: JWT tokens with dual auth (Magic Link + Password)

## Project Structure

```
TradePulse/
├── frontend/              # SvelteKit 5 application
│   ├── src/
│   │   ├── routes/       # File-based routing
│   │   └── lib/          # Components, stores, utilities
│   └── package.json
├── backend/               # Go API
│   ├── cmd/api/          # Application entry point
│   ├── internal/         # Private application code
│   └── migrations/       # Database migrations
├── docs/                  # Project documentation
│   ├── project.md        # Comprehensive project spec
│   └── api-spec.md       # API documentation
├── deployment/            # Deployment configurations
│   ├── nginx.conf.example
│   └── docker-compose.yml
└── .vscode/              # VSCode tasks and launch configs
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 20 or higher
- Access to PostgreSQL database at postgres1.drivenw.local
- External proxy configured for *.drivenw.com domains

### Initial Setup

1. **Clone the repository** (or create from template)

2. **Set up the database**

   The database schema will be automatically created on first run, or you can run migrations manually:

   ```bash
   psql -h postgres1.drivenw.local -p 5432 -U tradepulse -d tradepulse -f backend/migrations/001_initial_schema.up.sql
   ```

   Password: `drv13llc!`

3. **Configure Backend**

   ```bash
   cd backend
   cp .env.example .env
   ```

   Edit `.env` and set `JWT_SECRET` to a secure random string:
   ```bash
   # Generate a secure secret
   openssl rand -base64 32
   ```

4. **Configure Frontend**

   ```bash
   cd frontend
   cp .env.example .env
   ```

   The default API URL is already set to `https://api.tradepulse.drivenw.com:9000`

5. **Install Dependencies**

   Backend:
   ```bash
   cd backend
   go mod download
   ```

   Frontend:
   ```bash
   cd frontend
   npm install
   ```

### Running the Application

#### Using VSCode Tasks (Recommended)

1. Open the project in VSCode
2. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
3. Type "Run Task" and select "Tasks: Run Task"
4. Choose "Start All (Frontend + Backend)"

This will start both the backend and frontend servers in parallel.

#### Manual Start

**Backend:**
```bash
cd backend
go run cmd/api/main.go
```

The API will be available at: `https://api.tradepulse.drivenw.com:9000`

**Frontend:**
```bash
cd frontend
npm run dev
```

The frontend will be available at: `https://tradepulse.drivenw.com:4000`

### Accessing the Application

- **Frontend**: https://tradepulse.drivenw.com:4000
- **API**: https://api.tradepulse.drivenw.com:9000
- **API Health Check**: https://api.tradepulse.drivenw.com:9000/health

## Development

### Database Migrations

**Auto-Migration:** Migrations run automatically on backend startup.

Migrations are located in `backend/migrations/`:
- `001_initial_schema.up.sql` - Users, trades, tags, journal tables
- `002_add_password_auth.up.sql` - Password authentication

To run migrations manually (VSCode Task or command line):

```bash
# Apply migration
psql -h postgres1.drivenw.local -U tradepulse -d tradepulse -f backend/migrations/001_initial_schema.up.sql

# Rollback migration
psql -h postgres1.drivenw.local -U tradepulse -d tradepulse -f backend/migrations/001_initial_schema.down.sql
```

### Available VSCode Tasks

- **Start Backend** - Run Go API server
- **Start Frontend** - Run SvelteKit dev server
- **Start All** - Run both frontend and backend
- **Build Backend** - Compile Go binary
- **Build Frontend** - Build production frontend
- **Install Backend Dependencies** - Download Go modules
- **Install Frontend Dependencies** - Install npm packages
- **Run Database Migrations** - Apply initial schema migration
- **Run Password Migration** - Apply password auth migration

### Email Configuration

TradePulse supports multiple email providers for sending magic links:

1. **SMTP** (Default) - Supports trusted relay with blank credentials
2. **Microsoft Graph API** - For M365/Azure AD integration
3. **Google API** - For Gmail/Workspace integration

Configure the provider in `backend/.env` by setting `EMAIL_PROVIDER` and the corresponding credentials.

### Project Documentation

**Core Docs:**
- `/docs/project.md` - Complete project specification and architecture
- `/docs/authentication.md` - Authentication system guide (Magic Link + Password)
- `/docs/csv-import.md` - CSV import guide for DAS Trader and PropReports

**Backend:**
- `/docs/backend/implementation-status.md` - Current implementation status
- `/docs/backend/trades-api.md` - Trades API documentation
- `/docs/api-spec.md` - Full API specification (includes planned features)

**Frontend:**
- `/docs/frontend/components.md` - Component library documentation
- `/docs/frontend/design-system.md` - Design system and styling guide
- `/.claude/instructions.md` - Development guidelines and best practices

## API Documentation

The API follows RESTful conventions with consistent response formats.

**Base URL (Development)**: `https://api.tradepulse.drivenw.com:9000`

### Response Format

```json
{
  "success": true,
  "data": { ... }
}
```

### Key Endpoints

- `POST /api/auth/request-magic-link` - Request magic link
- `GET /api/auth/verify` - Verify magic link and get JWT
- `GET /api/trades` - List trades
- `POST /api/trades` - Create trade
- `GET /api/journal` - List journal entries
- `POST /api/journal` - Create journal entry
- `GET /api/metrics/summary` - Get performance metrics

See `/docs/api-spec.md` for complete API documentation.

## Deployment

### Development Environment

- Frontend runs on port 4000
- Backend API runs on port 9000
- External proxy handles *.drivenw.com routing
- Database at postgres1.drivenw.local:5432

### Production Deployment

For production deployment:

1. Use standard ports (80/443)
2. Configure SSL/TLS certificates
3. Set `ENVIRONMENT=production` in backend `.env`
4. Use strong `JWT_SECRET`
5. Enable database connection pooling
6. Configure file storage (S3 or similar for attachments)
7. Set up automated backups
8. Implement rate limiting

Reference nginx configuration is available in `deployment/nginx.conf.example`

## Security

- Passwordless authentication via magic links
- JWT tokens for session management
- Parameterized SQL queries to prevent injection
- File upload validation (type, size)
- CORS restricted to tradepulse.drivenw.com domain
- Secure password handling (bcrypt if passwords added)

## Contributing

This is a private project. For questions or issues, contact the project maintainers.

## Database Schema

Key tables:

- `users` - User accounts
- `magic_links` - Email authentication tokens
- `trades` - Trading history
- `journal_entries` - Trade journals with emotional state
- `attachments` - Screenshots and voice notes
- `tags` - Trade categorization

See `backend/migrations/001_initial_schema.up.sql` for complete schema.

## Future Enhancements

- Broker API integrations (TD Ameritrade, Interactive Brokers)
- Real-time notifications
- Mobile app
- Advanced analytics and pattern recognition
- Social features (anonymous sharing, mentorship)
- Custom metric calculations

## License

Proprietary - All rights reserved

## Support

For support or questions about TradePulse, contact the development team.

---

**Built with**: Go • SvelteKit • PostgreSQL • Love for Trading

# Claude Development Instructions for TradePulse

## Project Overview

TradePulse is a trading metrics and advanced journaling platform built with:
- **Frontend**: SvelteKit 5 + TypeScript
- **Backend**: Go API
- **Database**: PostgreSQL at postgres1.drivenw.local:5432

## Development Environment

### Domain & Port Configuration
- **Frontend**: https://tradepulse.drivenw.com (external) → port 4000 (internal)
- **API**: https://api.tradepulse.drivenw.com (external) → port 9000 (internal)
- **External Proxy**: Routes HTTPS (443) to internal non-standard ports
- **No localhost development** - all access via drivenw.com domains
- **No port specification needed** in API URLs or CORS configuration

### Database Connection
- Host: postgres1.drivenw.local
- Port: 5432
- Database: tradepulse
- User: tradepulse
- Password: drv13llc!
- SSL Mode: disable (internal network)

### File Locations
- Frontend: `/frontend`
- Backend: `/backend`
- Documentation: `/docs`
- Database Migrations: `/backend/migrations`

### Documentation Structure
- `/docs/architecture.md` - System architecture overview with diagrams
- `/docs/project.md` - Complete project specification
- `/docs/api-spec.md` - API endpoint reference
- `/docs/websocket-notifications.md` - WebSocket system guide
- `/docs/backend/` - Backend development guides
  - `getting-started.md` - Setup and running
  - `structure.md` - Code organization
  - `database.md` - Schema and migrations
  - `api-patterns.md` - Handler patterns
  - `notifications.md` - Publishing notifications
  - `email.md` - Email provider setup
- `/docs/frontend/` - Frontend development guides
  - `getting-started.md` - Setup and running
  - `README.md` - Quick reference

## Code Style & Standards

### Go Backend

**Project Structure:**
```
backend/
├── cmd/api/main.go           # Entry point
├── internal/                 # Private application code
│   ├── auth/                 # Authentication logic
│   ├── handlers/             # HTTP handlers
│   ├── models/               # Data models
│   ├── database/             # DB queries & connection
│   ├── email/                # Email providers
│   └── middleware/           # HTTP middleware
├── migrations/               # SQL migrations
└── go.mod
```

**Naming Conventions:**
- Package names: lowercase, single word (e.g., `auth`, `handlers`)
- File names: snake_case (e.g., `magic_link.go`, `jwt_token.go`)
- Exported functions: PascalCase (e.g., `GenerateMagicLink`)
- Unexported functions: camelCase (e.g., `validateEmail`)
- Constants: UPPER_SNAKE_CASE or PascalCase for exported

**Best Practices:**
- Use Go 1.21+ features (slices, maps packages)
- Prefer standard library over external dependencies
- Use context.Context for request-scoped values and cancellation
- Always close database connections and file handles with defer
- Use structured logging (slog package)
- Parameterized queries ONLY (prevent SQL injection)
- Error wrapping with fmt.Errorf and %w

**Error Handling:**
```go
if err != nil {
    return fmt.Errorf("failed to generate magic link: %w", err)
}
```

**Database Queries:**
```go
// Good - parameterized
row := db.QueryRow("SELECT id FROM users WHERE email = $1", email)

// Bad - vulnerable to SQL injection
row := db.QueryRow(fmt.Sprintf("SELECT id FROM users WHERE email = '%s'", email))
```

**API Response Format:**
Always use consistent JSON response structure:
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
}
```

### SvelteKit Frontend

**File Structure:**
```
frontend/src/
├── routes/                   # File-based routing
│   ├── +page.svelte         # Landing page
│   ├── +layout.svelte       # Root layout
│   ├── auth/
│   └── app/                 # Authenticated routes
├── lib/
│   ├── components/          # Reusable components
│   │   ├── ui/             # Skeleton UI wrappers
│   │   ├── trading/        # Trade-specific components
│   │   ├── journal/        # Journal components
│   │   ├── notifications/  # Notification UI
│   │   └── layout/         # Layout components
│   ├── stores/              # Svelte stores
│   ├── api/                 # API client
│   └── utils/               # Helper functions
```

**Component System:**
- **Base UI**: Skeleton UI library for rapid development
- **Custom Wrappers**: Thin wrappers in `components/ui/` for consistency
- **Domain Components**: Trading and journal-specific components
- **Dark Mode**: Built-in via Skeleton UI theme system

**Naming Conventions:**
- Components: PascalCase.svelte (e.g., `JournalEditor.svelte`)
- Utilities: camelCase.ts (e.g., `formatCurrency.ts`)
- Stores: camelCase.ts (e.g., `authStore.ts`)
- Types: PascalCase (e.g., `type Trade`, `interface User`)

**Component Structure:**
```svelte
<script lang="ts">
  // Imports
  import { onMount } from 'svelte';

  // Props
  export let trade: Trade;

  // State
  let isLoading = $state(false);

  // Derived state
  let pnlColor = $derived(trade.pnl >= 0 ? 'green' : 'red');

  // Functions
  function handleSubmit() {
    // ...
  }

  // Lifecycle
  onMount(() => {
    // ...
  });
</script>

<!-- Template -->
<div>
  {#if isLoading}
    <p>Loading...</p>
  {:else}
    <p style="color: {pnlColor}">{trade.pnl}</p>
  {/if}
</div>

<!-- Styles -->
<style>
  /* Component styles */
</style>
```

**API Client Pattern:**
```typescript
// lib/api/client.ts
const BASE_URL = import.meta.env.PUBLIC_API_URL;

export async function apiRequest<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const token = getAuthToken();

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options?.headers
    }
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.statusText}`);
  }

  return response.json();
}
```

**TypeScript:**
- Strict mode enabled
- No `any` types (use `unknown` if type is truly unknown)
- Define interfaces for all API responses
- Use type guards for runtime type checking

## Feature Development Guidelines

### When Adding New Features

1. **Plan First**: Review existing code structure in `/docs/project.md` and `/docs/api-spec.md`
2. **Database Changes**: Create migration files in `/backend/migrations/`
3. **Backend**:
   - Add model in `/backend/internal/models/`
   - Add handler in `/backend/internal/handlers/`
   - Register route in `main.go`
   - Update API documentation
4. **Frontend**:
   - Create/update components in `/frontend/src/lib/components/`
   - Add routes in `/frontend/src/routes/`
   - Update types in relevant files
5. **Testing**: Test both API endpoints and UI flows

### Security Checklist

Before implementing any feature, ensure:
- [ ] All SQL queries use parameterization
- [ ] User input is validated on backend
- [ ] Authentication middleware applied to protected routes
- [ ] File uploads validate type and size
- [ ] CORS configured for drivenw.com domain only
- [ ] Sensitive data not logged
- [ ] JWT tokens have reasonable expiry
- [ ] Magic links expire after 15 minutes

### Email System

The email system supports multiple providers via abstraction:

**Configuration:**
Set `EMAIL_PROVIDER` environment variable to: `smtp`, `graph`, or `gmail`

**SMTP Notes:**
- Supports blank username/password for trusted relay scenarios
- Set `SMTP_USER=""` and `SMTP_PASS=""` in .env for trusted relay

**Implementation Pattern:**
```go
type EmailProvider interface {
    Send(to, subject, body string) error
    SendHTML(to, subject, htmlBody string) error
}
```

When working with email:
- Use the abstraction, don't hardcode specific providers
- Handle errors gracefully (email failures shouldn't crash app)
- Log email sending for debugging
- Include unsubscribe link in production

## Common Tasks

### Adding a New API Endpoint

1. Define route in `backend/cmd/api/main.go`:
```go
r.Post("/api/endpoint", handlers.HandleEndpoint)
```

2. Create handler in `backend/internal/handlers/`:
```go
func HandleEndpoint(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

3. Update `/docs/api-spec.md` with endpoint documentation

### Adding a New Database Table

1. Create migration files:
```sql
-- backend/migrations/XXX_table_name.up.sql
CREATE TABLE table_name (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- backend/migrations/XXX_table_name.down.sql
DROP TABLE IF EXISTS table_name;
```

2. Add model in `backend/internal/models/table_name.go`

3. Update schema documentation in `/docs/project.md`

### Adding a New Frontend Route

1. Create route file: `frontend/src/routes/path/+page.svelte`
2. Create load function if needed: `frontend/src/routes/path/+page.ts`
3. Add navigation link in appropriate layout
4. Update types if new data structures introduced

### Handling File Uploads

**Backend:**
```go
func HandleUpload(w http.ResponseWriter, r *http.Request) {
    // Limit size
    r.ParseMultipartForm(10 << 20) // 10MB

    file, header, err := r.FormFile("file")
    if err != nil {
        // handle error
    }
    defer file.Close()

    // Validate MIME type
    // Save to disk or cloud storage
}
```

**Frontend:**
```svelte
<input type="file" accept="image/*" on:change={handleFileChange} />
```

## Database Patterns

### Queries
- Use prepared statements or query builders
- Always use `$1, $2, $3` placeholders (PostgreSQL)
- Return specific columns, avoid `SELECT *`
- Use transactions for multi-step operations

### Migrations
- Never modify existing migrations
- Always create new migration for schema changes
- Test both up and down migrations
- Keep migrations small and focused

## Environment Variables

### Backend Required Variables
```
PORT=9000
ENVIRONMENT=development
DB_HOST=postgres1.drivenw.local
DB_PORT=5432
DB_NAME=tradepulse
DB_USER=tradepulse
DB_PASSWORD=drv13llc!
JWT_SECRET=<generate-strong-secret>
EMAIL_PROVIDER=smtp
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASS=
SMTP_FROM=support@drivenw.com
ALLOWED_ORIGINS=https://tradepulse.drivenw.com
MAGIC_LINK_BASE_URL=https://tradepulse.drivenw.com
```

### Frontend Required Variables
```
PUBLIC_API_URL=https://api.tradepulse.drivenw.com
```

## Deployment Notes

- Backend runs on internal port 9000
- Frontend runs on internal port 4000
- External proxy routes HTTPS (443) to internal ports
- Public URLs: https://tradepulse.drivenw.com and https://api.tradepulse.drivenw.com
- No direct localhost access in development
- No port numbers in URLs, CORS, or API configuration
- PostgreSQL accessed via postgres1.drivenw.local:5432

## Git Workflow

1. Create feature branches from main
2. Meaningful commit messages: "Add journal entry voice recording feature"
3. Keep commits atomic and focused
4. Reference issue numbers if applicable

## Documentation Updates

When making significant changes, update:
- `/docs/project.md` - Overall architecture and features
- `/docs/api-spec.md` - API endpoints and contracts
- Code comments for complex logic
- This file if development workflow changes

## Troubleshooting

### Common Issues

**Database Connection Failed:**
- Verify postgres1.drivenw.local is accessible
- Check credentials: tradepulse / drv13llc!
- Ensure database `tradepulse` exists

**CORS Errors:**
- Verify `ALLOWED_ORIGINS` is set to `https://tradepulse.drivenw.com` (no port)
- Check frontend `.env` has `PUBLIC_API_URL=https://api.tradepulse.drivenw.com` (no port)
- Ensure Vite config includes `tradepulse.drivenw.com` in `allowedHosts`

**Magic Link Not Working:**
- Check email provider configuration
- Verify SMTP settings or provider credentials
- Check token expiry (15 minutes)

**Frontend Build Errors:**
- Clear `.svelte-kit` directory
- Delete `node_modules` and reinstall
- Check Node.js version (20+)

## Performance Considerations

- Use database indexes on frequently queried columns
- Implement pagination for large datasets (trades, journal entries)
- Lazy load images and attachments
- Use connection pooling for database
- Cache static assets
- Compress API responses

## Testing Strategy

### Backend
- Unit tests for business logic
- Integration tests for handlers
- Database tests with test database
- Email tests with mock provider

### Frontend
- Component tests with Vitest
- E2E tests for critical flows (auth, trade import, journal creation)
- Accessibility testing

## Future Considerations

As noted in `/docs/project.md`, future enhancements may include:
- Broker API integrations
- Real-time features (WebSocket)
- Mobile app
- Advanced analytics with AI

Keep architecture flexible to accommodate these features.

## Questions?

Refer to:
- `/docs/project.md` - Complete project specification
- `/docs/api-spec.md` - API documentation
- Code comments in existing implementations

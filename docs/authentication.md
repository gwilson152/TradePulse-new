# Authentication System

**Last Updated:** November 17, 2025

## Overview

TradePulse implements a dual authentication system allowing users to sign in via:
1. **Magic Link** (passwordless email verification)
2. **Email/Password** (traditional credentials)

The platform also supports a modern signup flow with plan selection, where users choose their subscription tier before creating their account. During the Beta period, all plans are free.

## User Flow

### First-Time User (Signup with Plan Selection)
1. Choose plan (Starter $2.99, Pro $9.99, Premium $14.99) - **Currently Free during Beta**
2. Enter email → Receive magic link via email
3. Click magic link → Authenticated with JWT token and plan activated
4. **Optional:** Set password in Settings → Account → Security
5. Next login: Choose Magic Link OR Password

### Alternative: Magic Link Without Plan Selection
1. Enter email → Receive magic link via email (defaults to Starter plan)
2. Click magic link → Authenticated with JWT token
3. **Optional:** Set password in Settings → Account → Security
4. Next login: Choose Magic Link OR Password

### Returning User (with password)
- Enter email + password → Instant authentication
- Can still use magic link if preferred

### Returning User (without password)
- Must use magic link
- Clear error message if password login attempted

## Backend Implementation

### Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/auth/signup` | POST | Public | Signup with plan selection |
| `/api/auth/request-magic-link` | POST | Public | Request magic link email |
| `/api/auth/verify` | GET | Public | Verify magic link token |
| `/api/auth/login` | POST | Public | Email/password login |
| `/api/auth/me` | GET | Protected | Get current user |
| `/api/auth/set-password` | POST | Protected | Set/update password |
| `/api/auth/logout` | POST | Protected | Logout (client-side) |

### Security Features

- **Password Hashing:** bcrypt with cost factor 10
- **Password Requirements:** Minimum 8 characters
- **JWT Tokens:** Configurable expiry (default 24h)
- **Magic Links:** One-time use, 15-minute expiry
- **Password Hash:** Never sent to client (excluded from User model JSON)
- **Race Condition:** Duplicate user creation handled gracefully

### Database Schema

```sql
-- Migration 002: Password Authentication
ALTER TABLE users ADD COLUMN password_hash TEXT;
CREATE INDEX idx_users_email_password ON users(email) WHERE password_hash IS NOT NULL;

-- Migration 003: User Plans
ALTER TABLE users ADD COLUMN plan_type VARCHAR(50) DEFAULT 'starter';
ALTER TABLE users ADD COLUMN plan_status VARCHAR(50) DEFAULT 'beta_free';
ALTER TABLE users ADD COLUMN plan_selected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
CREATE INDEX idx_users_plan_type ON users(plan_type);
```

**User Model:**
- `id` (UUID)
- `email` (string, unique)
- `password_hash` (string, internal only)
- `has_password` (boolean, sent to client)
- `plan_type` (string: 'starter', 'pro', 'premium')
- `plan_status` (string: 'beta_free', 'active', 'cancelled', 'trial', 'expired')
- `plan_selected_at` (timestamp)
- `created_at` (timestamp)
- `last_login` (timestamp)
- `preferences` (JSONB)

### Plan System

**Available Plans:**
- **Starter** - $2.99/month (Free during Beta)
- **Pro** - $9.99/month (Free during Beta)
- **Premium** - $14.99/month (Free during Beta)

**Beta Status:**
All users during the Beta period have `plan_status = 'beta_free'`, meaning they get full access to their selected plan features without payment. This allows users to experience the platform before the official launch.

**Plan Fields:**
- `plan_type`: The selected subscription tier ('starter', 'pro', 'premium')
- `plan_status`: Current plan status ('beta_free', 'active', 'cancelled', 'trial', 'expired')
- `plan_selected_at`: Timestamp when the user selected their plan

### Auto-Migrations

Migrations run automatically on server startup via `db.RunMigrations()`:
- Migration 001: Initial schema (users, trades, tags, journal)
- Migration 002: Password authentication
- Migration 003: User plan fields

## Frontend Implementation

### Login Page (`/auth/login`)

**Dual-Mode UI:**
- Toggle between "Password" and "Magic Link" modes
- Default mode: Password (faster for returning users)
- Real-time validation and error messages
- Loading states and accessibility

**Features:**
- Email input (both modes)
- Password input (password mode only)
- Mode toggle button
- Clear error messages for edge cases

### User Menu Component

**Location:** Top-right navigation bar

**Features:**
- User avatar with initials
- Dropdown menu on click
- User email display
- Settings navigation
- Logout functionality
- Click-outside & Escape key to close

### Password Management

**Location:** Settings → Account → Security

**Features:**
- Shows password status (set or not set)
- "Set Password" or "Change Password" button
- Password strength indicator
- Real-time validation
- Visual requirements checklist

**Password Requirements:**
- ✅ Minimum 8 characters
- ✅ Mix of uppercase and lowercase (recommended)
- ✅ At least one number (recommended)

### API Client

```typescript
// Password authentication
await apiClient.loginWithPassword(email, password);

// Set password (requires authentication)
await apiClient.setPassword(password);

// Get current user
const user = await apiClient.getCurrentUser();

// Magic link authentication
await apiClient.requestMagicLink(email);
await apiClient.verifyMagicLink(token);

// Logout
await apiClient.logout();

// Decode JWT token (fallback for user info)
const payload = apiClient.getTokenPayload();
```

### User Store

```typescript
import { userStore } from '$lib/stores/user';

// Set user
userStore.setUser(user);

// Clear user
userStore.clearUser();

// Helper functions
const initials = getUserInitials(user);  // "JD"
const displayName = getDisplayName(user); // "John"
```

## Configuration

### Environment Variables

```bash
# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h

# Magic Link
MAGIC_LINK_EXPIRY=15m
MAGIC_LINK_BASE_URL=https://tradepulse.drivenw.com

# Email (stubbed - logs to console)
EMAIL_PROVIDER=smtp
SMTP_HOST=your-smtp-host
SMTP_PORT=587
```

## Troubleshooting

### 401 Unauthorized Errors

**Symptoms:** API calls to protected endpoints return 401

**Common Causes:**
1. **Token not stored:** Check if JWT token is in localStorage (`auth_token`)
2. **Backend not restarted:** After code changes, restart the backend server
3. **Wrong context key:** Ensure handlers use `middleware.GetUserID(r)` not `r.Context().Value("user_id")`

**Solution:**
```bash
# Check browser console for token
localStorage.getItem('auth_token')

# Restart backend
cd backend && go run cmd/api/main.go
```

### Password Not Saving

**Symptoms:** Password modal shows success but user still has no password

**Causes:**
- Migration 002 not applied
- Database connection issue

**Solution:**
```bash
# Check if password_hash column exists
psql -h postgres1.drivenw.local -U tradepulse -d tradepulse -c "\d users"

# Run migration if needed (auto-runs on startup)
cd backend && go run cmd/api/main.go
```

### Magic Link Not Received

**Symptoms:** No email arrives after requesting magic link

**Cause:** Email sending is stubbed

**Solution:**
Check backend console logs for the magic link URL:
```
{"level":"INFO","msg":"Magic link generated","token":"...","link":"https://..."}
```

## Known Limitations

- 📧 **Email Sending:** Currently stubbed - magic links logged to console
- 🔄 **Token Refresh:** Endpoint planned but not implemented
- 🔐 **2FA:** Not currently supported
- 📱 **Session Management:** Basic JWT only (no refresh tokens)

## Future Enhancements

- Email provider integration (SMTP/Microsoft Graph/Gmail)
- JWT refresh token mechanism
- Two-factor authentication (TOTP)
- Session management dashboard
- Password reset via email
- Account recovery options

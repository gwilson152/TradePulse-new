# Authentication System

**Last Updated:** November 16, 2025

## Overview

TradePulse implements a dual authentication system allowing users to sign in via:
1. **Magic Link** (passwordless email verification)
2. **Email/Password** (traditional credentials)

## User Flow

### First-Time User
1. Enter email → Receive magic link via email
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
```

**User Model:**
- `id` (UUID)
- `email` (string, unique)
- `password_hash` (string, internal only)
- `has_password` (boolean, sent to client)
- `created_at` (timestamp)
- `last_login` (timestamp)

### Auto-Migrations

Migrations run automatically on server startup via `db.RunMigrations()`:
- Migration 001: Initial schema (users, trades, tags, journal)
- Migration 002: Password authentication

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

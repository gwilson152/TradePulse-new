# Authentication

TradePulse uses passwordless authentication with magic links sent via email.

## How It Works

1. **User requests magic link** - POST email to `/api/auth/request-magic-link`
2. **System generates token** - 64-character secure hex token with 15-minute expiry
3. **Token stored in database** - One-time use token linked to user
4. **Magic link sent** - Email with link to `/auth/verify?token=...` (currently logged to console)
5. **User clicks link** - Browser visits verification URL
6. **Token verified** - Backend checks token validity and marks as used
7. **JWT issued** - 24-hour JWT token returned to frontend
8. **User authenticated** - JWT used for subsequent API requests

## API Endpoints

### Request Magic Link

**POST** `/api/auth/request-magic-link`

Request body:
```json
{
  "email": "user@example.com"
}
```

Success response (200):
```json
{
  "success": true,
  "data": {
    "message": "Magic link sent to your email"
  }
}
```

Error response (400):
```json
{
  "success": false,
  "error": {
    "code": "INVALID_EMAIL",
    "message": "Email is required"
  }
}
```

**Behavior**:
- Creates new user if email doesn't exist
- Generates secure 64-character hex token
- Stores token with 15-minute expiry
- Deletes any existing unused tokens for the user
- Logs magic link to console (email sending TODO)

### Verify Magic Link

**GET** `/api/auth/verify?token=abc123...`

Success response (200):
```json
{
  "success": true,
  "data": {
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "created_at": "2025-01-10T14:00:00Z",
      "last_login": "2025-01-10T14:05:00Z"
    }
  }
}
```

Error responses:
- `INVALID_TOKEN` (401) - Token not found, expired, or already used
- `USER_FETCH_ERROR` (500) - Failed to retrieve user data

**Behavior**:
- Validates token exists and hasn't been used
- Checks token hasn't expired (15 minutes)
- Marks token as used (one-time use)
- Updates user's last_login timestamp
- Generates JWT with 24-hour expiry
- Returns JWT and user data

## JWT Token

### Claims

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "exp": 1704902400,
  "iat": 1704816000
}
```

### Using JWT

Include JWT in Authorization header for protected endpoints:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

The JWT middleware validates the token and adds user_id to request context.

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE,
    preferences JSONB DEFAULT '{}'::jsonb
);
```

### Magic Links Table

```sql
CREATE TABLE magic_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Implementation Details

### Token Generation

```go
tokenBytes := make([]byte, 32)
rand.Read(tokenBytes)
token := hex.EncodeToString(tokenBytes) // 64 characters
```

### JWT Signing

```go
claims := jwt.MapClaims{
    "user_id": userID.String(),
    "email":   user.Email,
    "exp":     time.Now().Add(24 * time.Hour).Unix(),
    "iat":     time.Now().Unix(),
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signedToken, _ := token.SignedString([]byte(jwtSecret))
```

### Token Verification

1. Query database for token
2. Check if `used_at IS NULL`
3. Check if `expires_at > NOW()`
4. Update `used_at = NOW()`
5. Return user_id

## Security Features

- **Secure random tokens** - 256-bit cryptographically secure random tokens
- **Time-limited** - Tokens expire after 15 minutes
- **One-time use** - Tokens marked as used after verification
- **Auto-cleanup** - Old unused tokens deleted when new one requested
- **JWT expiry** - JWT tokens expire after 24 hours
- **HTTPS only** - All authentication endpoints require HTTPS in production

## Testing Authentication

### Development Testing

Since email sending is not yet implemented, magic links are logged to the console:

1. Start the backend server
2. Make a request:
```bash
curl -X POST https://api.tradepulse.drivenw.com/api/auth/request-magic-link \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
```
3. Check backend logs for the magic link URL
4. Visit the URL to get your JWT

### Frontend Testing

1. Go to https://tradepulse.drivenw.com/auth/login
2. Enter your email
3. Check backend console for magic link
4. Copy and visit the link
5. You'll be redirected to dashboard with authentication

## Future Enhancements

- [ ] Implement email sending (SMTP/Graph/Gmail)
- [ ] Add rate limiting on magic link requests
- [ ] Implement JWT refresh tokens
- [ ] Add session management
- [ ] Email verification status
- [ ] Two-factor authentication option

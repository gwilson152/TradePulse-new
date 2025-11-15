# Frontend Getting Started

## Prerequisites

- Node.js 20 or higher
- npm or pnpm

## Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env`:

```bash
PUBLIC_API_URL=https://api.tradepulse.drivenw.com
```

### 3. Start Dev Server

```bash
npm run dev
```

The app will start on port 4000 and be accessible at:
`https://tradepulse.drivenw.com`

## Verify Installation

1. Open https://tradepulse.drivenw.com
2. Should see landing page
3. Check browser console for errors

## Development Commands

### Run Dev Server
```bash
npm run dev
```

### Build for Production
```bash
npm run build
```

### Preview Production Build
```bash
npm run preview
```

### Type Check
```bash
npm run check
```

### Format Code
```bash
npm run format  # if configured
```

## Project Structure

```
frontend/
├── src/
│   ├── routes/              # Pages (file-based routing)
│   │   ├── +page.svelte    # Landing page
│   │   ├── +layout.svelte  # Root layout
│   │   └── app/            # Authenticated routes
│   ├── lib/                 # Reusable code
│   │   ├── components/     # UI components
│   │   ├── stores/         # State management
│   │   ├── api/            # API communication
│   │   └── utils/          # Helper functions
│   ├── app.html            # HTML template
│   └── app.css             # Global styles
├── static/                  # Public assets
├── .env                     # Environment variables
├── svelte.config.js         # SvelteKit config
├── vite.config.ts           # Vite config
├── tailwind.config.js       # Tailwind CSS config
└── package.json             # Dependencies
```

## Common Issues

### Port Already in Use

**Error**: `Port 4000 is already in use`

**Solution**:
```bash
# Find process
lsof -i :4000  # Mac/Linux
netstat -ano | findstr :4000  # Windows

# Kill process
kill -9 <PID>
```

### Host Not Allowed

**Error**: `Blocked request. This host ("tradepulse.drivenw.com") is not allowed.`

**Solution**: Verify `vite.config.ts` includes `allowedHosts`:
```typescript
server: {
  allowedHosts: ['tradepulse.drivenw.com']
}
```

### API Connection Failed

**Error**: CORS or connection errors

**Solution**:
- Verify backend is running
- Check `PUBLIC_API_URL` in `.env`
- Verify external proxy is routing correctly

## Next Steps

- Read [Structure](structure.md) to understand code organization
- See [Components](components.md) for UI component usage
- Check [API Client](api-client.md) for backend communication
- Review [State Management](state-management.md) for stores pattern

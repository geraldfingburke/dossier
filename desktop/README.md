# Dossier Desktop Application

Electron wrapper for the Dossier news digest application, providing a native desktop experience.

## Architecture

The desktop application consists of two main components:

1. **Electron Wrapper** (this directory)

   - Native window management
   - System tray integration
   - Auto-update support
   - IPC bridge for secure frontend-backend communication

2. **Vue.js Frontend** (loaded from `../client`)

   - Existing web client reused without modification
   - Communicates with backend via GraphQL API

3. **Go Backend Service** (separate process)
   - Runs as system service or Docker container
   - Provides GraphQL API on `http://localhost:8080`

## Development Setup

### Prerequisites

- Node.js 18+
- Docker (for running the backend service)
- npm or yarn

### Installation

```bash
# Install dependencies
cd desktop
npm install

# Start the backend service (in separate terminal)
cd ..
docker-compose up

# Build the Vue.js client (in separate terminal)
cd client
npm install
npm run dev

# Run Electron in development mode
cd ../desktop
npm run dev
```

### Development Mode

In development mode (`npm run dev`):

- Electron loads the Vue.js app from Vite dev server (`http://localhost:5173`)
- Backend service expected at `http://localhost:8080` (via Docker)
- DevTools automatically opened
- Hot reload enabled for frontend changes

## Building for Production

### Build Client First

```bash
# Build the Vue.js client for production
cd client
npm run build

# This creates client/dist/ which gets bundled into Electron
```

### Build Desktop App

```bash
cd desktop

# Build for all platforms
npm run build

# Or build for specific platform
npm run build:win    # Windows
npm run build:mac    # macOS
npm run build:linux  # Linux
```

### Output

Built applications are in `desktop/dist/`:

- **Windows**: `.exe` installer and portable `.exe`
- **macOS**: `.dmg` installer and `.zip` archive
- **Linux**: `.AppImage`, `.deb`, and `.rpm` packages

## Project Structure

```
desktop/
├── src/
│   ├── main.js         # Electron main process
│   └── preload.js      # Context bridge (IPC)
├── build/
│   ├── icon.ico        # Windows icon
│   ├── icon.icns       # macOS icon
│   └── icon.png        # Linux icon
├── package.json        # Electron config & dependencies
└── README.md          # This file
```

## Features

### Window Management

- Minimize to system tray
- Remember window size/position (TODO)
- Native menus (File, Edit, View, Help)

### System Integration

- System tray icon with menu
- Native notifications
- Auto-update support (GitHub releases)
- Platform-specific installers

### Security

- Context isolation enabled
- Node integration disabled
- Secure IPC bridge via `contextBridge`
- External links open in default browser

## Configuration

### Backend Service URL

Default: `http://localhost:8080`

To change, edit `servicePort` in `src/main.js`:

```javascript
const servicePort = 8080;
```

### Auto-Update

Configure in `package.json` under `build.publish`:

```json
"publish": {
  "provider": "github",
  "owner": "geraldfingburke",
  "repo": "dossier"
}
```

Updates are checked:

- On application startup (after 3 second delay)
- Via "Check for Updates" in tray menu
- Automatically downloaded and installed

## API Exposed to Frontend

The preload script exposes these APIs to the Vue.js app via `window.electronAPI`:

```javascript
// Service communication
await window.electronAPI.getServiceUrl();
await window.electronAPI.checkServiceHealth();

// External links
await window.electronAPI.openExternal(url);

// App info
await window.electronAPI.getAppVersion();

// Window management
await window.electronAPI.minimizeToTray();

// Notifications
await window.electronAPI.showNotification({ title, body });

// Navigation listener
window.electronAPI.onNavigate((route) => {
  /* ... */
});
```

## Platform Information

Available via `window.platform`:

```javascript
window.platform.os; // 'win32', 'darwin', or 'linux'
window.platform.isWindows; // boolean
window.platform.isMac; // boolean
window.platform.isLinux; // boolean
```

## TODO

### Phase 1 (Current)

- [x] Basic Electron wrapper
- [x] Window management
- [x] System tray integration
- [x] IPC bridge
- [x] Development mode support
- [ ] App icons (placeholder currently)

### Phase 2 (Next)

- [ ] Remember window size/position
- [ ] Keyboard shortcuts
- [ ] Deep linking support
- [ ] Native file dialogs (export dossiers)

### Phase 3 (Future)

- [ ] Bundle Go service with app
- [ ] SQLite support (vs PostgreSQL)
- [ ] Offline mode
- [ ] System service installer (Windows/macOS/Linux)

## Backend Service Integration

Currently expects backend at `http://localhost:8080` via Docker:

```bash
docker-compose up
```

**Future**: Bundle Go service binary with Electron app and manage as child process or system service.

## Troubleshooting

### "Service Not Running" error

**Cause**: Backend service not accessible at `http://localhost:8080`

**Solution**:

```bash
# Start backend via Docker
docker-compose up

# Or check if service is running
curl http://localhost:8080/health
```

### Blank window in development

**Cause**: Vite dev server not running

**Solution**:

```bash
cd client
npm run dev
```

### Build fails with "client/dist not found"

**Cause**: Vue.js client not built for production

**Solution**:

```bash
cd client
npm run build
```

## License

MIT

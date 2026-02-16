# Build & Development Guide

This guide covers building, testing, and developing the FileBrowser Quantum - ChainFS Integration Fork.

Note: all builds are done in powershell as development is happening on windows.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Building the Application](#building-the-application)
- [Running the Application](#running-the-application)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Building for Production](#building-for-production)
- [Additional Resources](#additional-resources)
- [Support](#support)

---

## Prerequisites

### Required Software

**Backend (Go):**
- Go 1.25.0 or higher
- Git

**Frontend (Vue.js):**
- Node.js 18+ (LTS recommended)
- npm 9+ or yarn

**Operating Systems:**
- Windows 10/11
- Linux (Ubuntu 20.04+, Debian, etc.)
- macOS 11+

### Optional Tools
- Docker (for containerized deployment)
- Azure CLI (for cloud deployment)

---

## Quick Start

**Windows (PowerShell):**

# 1. Build the frontend
```powershell
cd C:/git/filebrowser2/frontend
npm install
npm run build:windows
```

# 2. Build the backend
```powershell
cd C:/git/filebrowser2/backend
go build -o filebrowser.exe
```

# 3. Run the application
```powershell
cd C:/git/filebrowser2/backend
./filebrowser.exe -c config.dev.yaml
```

Open browser to: `http://localhost:8080`

---

## Building the Application

### Frontend Build

The frontend is a Vue.js 3 application that must be built before running the backend.

**Windows:**
```powershell
cd C:/git/filebrowser2/frontend

# Install dependencies
npm install

# Development build (with hot reload)
npm run dev

# Production build (Windows-specific)
npm run build:windows
```

**Note:** The build script uses Unix commands (`rm`, `cp`) on Linux/macOS. On Windows, use `npm run build:windows` which uses a PowerShell script, or run the build through Git Bash/WSL.

**Build Output:**
- Production build creates files in: `frontend/dist/`
- These files are embedded into the Go binary during backend build

**Frontend Development Server:**

Is launched by the backend. do not try to run seperately.

### Backend Build

The backend is written in Go and embeds the frontend assets.

```powershell
cd C:/git/filebrowser2/backend

# Development build
go build -o filebrowser.exe

# Production build (optimized)
go build -ldflags="-s -w" -o filebrowser.exe
```

**Build Flags:**
- `-ldflags="-s -w"` - Strip debug info and symbol table (smaller binary)
- `-o` - Specify output filename

**Binary Size:**
- Development: ~45-50 MB
- Production (stripped): ~35-40 MB

---

## Running the Application

### Configuration Files

The application uses YAML configuration files:

- `config.yaml` - Default configuration (uses DEV settings)
- `config.dev.yaml` - DEV environment (points to nansendev.azurewebsites.net)
- `config.uat.yaml` - UAT environment (points to nansenuat.azurewebsites.net)
- `config.prod.yaml` - PROD environment (points to nansenprod.azurewebsites.net)

### Starting the Server

**Using specific config:**
```powershell
cd C:/git/filebrowser2/backend

# DEV environment (for all testing use DEV config)
./filebrowser.exe -c config.dev.yaml

# UAT environment
./filebrowser.exe -c config.uat.yaml

# PROD environment
./filebrowser.exe -c config.prod.yaml
```

**Command-line flags:**
```powershell
./filebrowser.exe -h              # Show help
./filebrowser.exe -c              # Print default config
./filebrowser.exe version         # Show version info
```

### Default Ports

- **DEV/UAT:** `8080`
- **PROD (main config):** `80`

### Accessing the Application

After starting the server:
1. Open browser to `http://localhost:8080` (or configured port)
2. You should see the login page with **"ChainFS Login"** button
3. Click "ChainFS Login" to authenticate via Azure AD B2C

### Initial Admin User

**Password authentication (disabled by default):**
- Username: `admin`
- Password: `admin`

**ChainFS authentication (enabled by default):**
- Users are auto-created on first Azure AD B2C login
- Admin status determined by Azure AD roles/groups claim

---

## Testing

### Backend Tests

```powershell
cd C:/git/filebrowser2/backend

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./http/...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestChainFsLogin ./http/...
```

### Frontend Tests

```powershell
cd C:/git/filebrowser2/frontend

# Run unit tests
npm run test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch
```

### Database Inspection

```bash
# View database contents (requires sqlite3)
sqlite3 database.db

# List all users
sqlite> SELECT id, username, loginMethod FROM users;

# Check Azure token fields (encrypted)
sqlite> SELECT username, azureAccessToken, azureTokenExpiry FROM users;

# Exit
sqlite> .quit
```

---

## Troubleshooting

### Common Build Issues

**1. Frontend build fails**
```bash
# Clear node_modules and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

**2. Backend build fails with "template: pattern matches no files"**
```powershell
# Frontend not built yet
cd C:/git/filebrowser2/frontend
npm run build

# Then rebuild backend
cd C:/git/filebrowser2/backend
go build -o filebrowser.exe
```

**3. Go module issues**
```powershell
cd C:/git/filebrowser2/backend
go mod tidy
go mod download
go build -o filebrowser.exe
```

### Common Runtime Issues

**1. "Auth Methods: [password]" instead of "Auth Methods: [chainfs]"**

**Problem:** ChainFS auth not enabled in config

**Solution:**
```yaml
# config.dev.yaml
auth:
  methods:
    chainfs:
      enabled: true  # Must be true
    password:
      enabled: false # Must be false
```

**2. "Invalid redirect URI" error from Azure AD B2C**

**Problem:** Redirect URI not registered in Azure AD B2C

**Solution:**
- Register callback URL in Azure AD B2C application settings:
  ```
  http://localhost:8080/api/auth/chainfs/callback
  ```
- Ensure exact match (protocol, domain, port, path)

**3. "User does not exist" error**

**Problem:** `createUser: false` but user not in database

**Solution:**
```yaml
# config.dev.yaml
auth:
  methods:
    chainfs:
      createUser: true  # Enable auto-user creation
```

**4. Port already in use**

**Problem:** Another process using port 8080

**Solution:**
```powershell
netstat -ano | findstr :8080
```

**5. Database locked**

**Problem:** Multiple FileBrowser instances accessing same database

**Solution:**
- Stop all FileBrowser instances
- Delete `database.db-shm` and `database.db-wal` files
- Restart FileBrowser

### Viewing Logs

**Enable verbose logging:**
```yaml
# config.yaml
server:
  logging:
    - levels: "debug|info|warning|error"
```

**Log locations:**
- Console output (stdout/stderr)
- No default log files (logs to console only)

---

## Building for Production

### Complete Production Build

```powershell
# 1. Build frontend (production mode)
cd C:/git/filebrowser2/frontend
npm install --production
npm run build

# 2. Build backend (optimized)
cd C:/git/filebrowser2/backend
go build -ldflags="-s -w" -o filebrowser.exe

# 3. Prepare distribution
mkdir -p dist
copy filebrowser.exe dist/
copy config.prod.yaml dist/config.yaml

# 4. Create tarball
cd dist
tar -czf filebrowser-chainfs-v1.0.0.tar.gz *
```

### Useful Commands

```powershell
# Check Go version
go version

# Check Node version
node --version

# Check npm version
npm --version

# View FileBrowser version
./filebrowser.exe version

# Clean build artifacts
cd backend && go clean
cd frontend && rm -rf dist node_modules
```

---

## Additional Resources

- **Main Documentation:** [Claude.md](Claude.md)
- **Original FileBrowser:** https://github.com/gtsteffaniak/filebrowser
- **ChainFS API Docs:** Check Swagger endpoints at ChainFS API URLs
- **Vue.js 3 Docs:** https://vuejs.org/
- **Go Documentation:** https://golang.org/doc/

---

## Support

For upstream FileBrowser issues:
- Original repository: https://github.com/gtsteffaniak/filebrowser
- Documentation: https://github.com/gtsteffaniak/filebrowser/wiki


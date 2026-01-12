# Build & Development Guide

This guide covers building, testing, and developing the FileBrowser Quantum - ChainFS Integration Fork.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Building the Application](#building-the-application)
- [Running the Application](#running-the-application)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

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

**Linux/macOS:**
```bash
# 1. Clone the repository
git clone https://github.com/your-username/filebrowser-chainfs.git
cd filebrowser-chainfs

# 2. Build the frontend
cd frontend
npm install
npm run build

# 3. Build the backend
cd ../backend
go build -o filebrowser

# 4. Run the application
./filebrowser -c config.dev.yaml
```

**Windows (PowerShell):**
```powershell
# 1. Clone the repository
git clone https://github.com/your-username/filebrowser-chainfs.git
cd filebrowser-chainfs

# 2. Build the frontend
cd frontend
npm install
npm run build:windows

# 3. Build the backend
cd ..\backend
go build -o filebrowser.exe

# 4. Run the application
.\filebrowser.exe -c config.dev.yaml
```

Open browser to: `http://localhost:8080`

---

## Building the Application

### Frontend Build

The frontend is a Vue.js 3 application that must be built before running the backend.

**Linux/macOS:**
```bash
cd frontend

# Install dependencies
npm install

# Development build (with hot reload)
npm run dev

# Production build (optimized)
npm run build
```

**Windows:**
```powershell
cd frontend

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
```bash
npm run dev
# Runs on http://localhost:5173
# Proxies API requests to backend on http://localhost:8080
```

### Backend Build

The backend is written in Go and embeds the frontend assets.

```bash
cd backend

# Development build
go build -o filebrowser

# Production build (optimized)
go build -ldflags="-s -w" -o filebrowser

# Cross-compilation examples
# Windows
GOOS=windows GOARCH=amd64 go build -o filebrowser.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o filebrowser

# macOS
GOOS=darwin GOARCH=amd64 go build -o filebrowser
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
```bash
cd backend

# DEV environment
./filebrowser -c config.dev.yaml

# UAT environment
./filebrowser -c config.uat.yaml

# PROD environment
./filebrowser -c config.prod.yaml

# Default (uses config.yaml)
./filebrowser
```

**Command-line flags:**
```bash
./filebrowser -h              # Show help
./filebrowser -c              # Print default config
./filebrowser version         # Show version info
```

**Environment Variables:**
```bash
# Custom config path
export FILEBROWSER_CONFIG=/path/to/config.yaml
./filebrowser

# Generate config.yaml template
export FILEBROWSER_GENERATE_CONFIG=true
./filebrowser
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

## Development Workflow

### Backend Development

**1. Make code changes**
```bash
cd backend
# Edit .go files
```

**2. Rebuild**
```bash
go build -o filebrowser
```

**3. Run**
```bash
./filebrowser -c config.dev.yaml
```

**Hot reload (using tools):**
```bash
# Install air (Go hot reload tool)
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Frontend Development

**1. Start backend** (in one terminal):
```bash
cd backend
./filebrowser -c config.dev.yaml
```

**2. Start frontend dev server** (in another terminal):
```bash
cd frontend
npm run dev
```

**3. Access frontend dev server:**
- Frontend: `http://localhost:5173` (with hot reload)
- API requests automatically proxy to backend on `http://localhost:8080`

**Frontend file structure:**
```
frontend/
├── src/
│   ├── api/          # API client functions
│   ├── components/   # Vue components
│   ├── router/       # Vue Router config
│   ├── store/        # State management
│   ├── utils/        # Utility functions
│   └── views/        # Page components
├── public/           # Static assets
├── package.json      # Node dependencies
└── vite.config.ts    # Vite build config
```

### Making Changes to Authentication

**Backend changes:**
1. Modify `backend/http/chainfs.go` - Authentication handlers
2. Modify `backend/chainfs/client.go` - ChainFS API client
3. Rebuild: `go build -o filebrowser`

**Frontend changes:**
1. Modify `frontend/src/views/Login.vue` - Login UI
2. Hot reload automatically applies changes (if using `npm run dev`)
3. For production: `npm run build` then rebuild backend

**Configuration changes:**
1. Edit `backend/config.dev.yaml` (or other config files)
2. Restart backend: `./filebrowser -c config.dev.yaml`

---

## Testing

### Backend Tests

```bash
cd backend

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

```bash
cd frontend

# Run unit tests
npm run test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch
```

### Integration Testing

**Manual testing checklist:**

1. **ChainFS Login Flow:**
   - [ ] Click "ChainFS Login" button
   - [ ] Redirected to Azure AD B2C login page
   - [ ] Enter valid credentials
   - [ ] Redirected back to FileBrowser
   - [ ] Logged in successfully (JWT cookie set)
   - [ ] Can browse files

2. **User Creation:**
   - [ ] New user auto-created on first login (if `createUser: true`)
   - [ ] User stored in `database.db`
   - [ ] Azure tokens encrypted in database

3. **Admin Permissions:**
   - [ ] User with admin role claim has admin access
   - [ ] User without admin role has normal access

4. **Logout:**
   - [ ] Click logout
   - [ ] Redirected to Azure AD B2C logout
   - [ ] Session cleared
   - [ ] Cannot access files without re-login

5. **Token Persistence:**
   - [ ] Close browser
   - [ ] Reopen FileBrowser
   - [ ] Still logged in (session persists)

### Testing Against Different Environments

**DEV (Development):**
```bash
./filebrowser -c config.dev.yaml
# Test against: https://nansendev.azurewebsites.net
```

**UAT (User Acceptance Testing):**
```bash
./filebrowser -c config.uat.yaml
# Test against: https://nansenuat.azurewebsites.net
```

**PROD (Production):**
```bash
./filebrowser -c config.prod.yaml
# Test against: https://nansenprod.azurewebsites.net
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
```bash
# Frontend not built yet
cd frontend
npm run build

# Then rebuild backend
cd ../backend
go build -o filebrowser
```

**3. Go module issues**
```bash
cd backend
go mod tidy
go mod download
go build -o filebrowser
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
```bash
# Find process using port
# Linux/Mac:
lsof -i :8080

# Windows:
netstat -ano | findstr :8080

# Kill process or change port in config
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

**Debugging authentication:**
```bash
# Run with verbose output
./filebrowser -c config.dev.yaml 2>&1 | tee filebrowser.log
```

### Performance Issues

**1. Slow frontend loading**
- Use production build: `npm run build` (not `npm run dev`)
- Frontend assets are embedded and served from memory

**2. Slow backend startup**
- First run creates database and indexes
- Subsequent starts are faster
- Large file directories may take time to index

---

## Building for Production

### Complete Production Build

```bash
# 1. Build frontend (production mode)
cd frontend
npm install --production
npm run build

# 2. Build backend (optimized)
cd ../backend
go build -ldflags="-s -w" -o filebrowser

# 3. Prepare distribution
mkdir -p dist
cp filebrowser dist/
cp config.prod.yaml dist/config.yaml
cp Claude.md dist/
cp BUILD.md dist/

# 4. Create tarball
cd dist
tar -czf filebrowser-chainfs-v1.0.0.tar.gz *
```

### Docker Build

```bash
# Build Docker image
docker build -t filebrowser-chainfs:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -v /path/to/files:/srv \
  -v /path/to/database:/app/database.db \
  -e FILEBROWSER_CONFIG=/app/config.prod.yaml \
  filebrowser-chainfs:latest
```

### Azure Deployment

See **Priority 3: Azure Hosting** in `Claude.md` for deployment instructions.

---

## Development Tips

### Code Style

**Go:**
- Follow standard Go conventions
- Run `gofmt` before committing
- Use `golangci-lint` for linting

**Vue.js/TypeScript:**
- Follow Vue.js 3 Composition API style
- Use ESLint for linting
- Run `npm run lint` before committing

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make changes and commit
git add .
git commit -m "feat: add feature description"

# Push to remote
git push origin feature/your-feature-name

# Create pull request
```

### Commit Message Convention

```
feat: Add new feature
fix: Fix bug
docs: Update documentation
style: Code style changes
refactor: Code refactoring
test: Add tests
chore: Maintenance tasks
```

### Useful Commands

```bash
# Check Go version
go version

# Check Node version
node --version

# Check npm version
npm --version

# View FileBrowser version
./filebrowser version

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

For issues specific to this fork:
1. Check this BUILD.md guide
2. Review [Claude.md](Claude.md) for architecture details
3. Check ChainFS API source code: `C:\git\azure-blockchain-workbench-app\NasenAPI`

For upstream FileBrowser issues:
- Original repository: https://github.com/gtsteffaniak/filebrowser
- Documentation: https://github.com/gtsteffaniak/filebrowser/wiki

---

**Last Updated:** 2026-01-13
**Version:** 1.0.0
**Status:** Priority 1 (Authentication) Complete

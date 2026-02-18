# Build & Development Guide

This guide covers building, testing, and developing the FileBrowser Quantum - ChainFS Integration Fork.

## Table of Contents
- [Prerequisites](#prerequisites)
  - [Required Software](#required-software)
  - [Optional Tools](#optional-tools)
- [Quick Start (DEV Build)](#quick-start-dev-build)
  - [Frontend](#frontend)
  - [Backend](#backend)
  - [Run](#run)
- [Running the Application](#running-the-application)
  - [Configuration Files](#configuration-files)
- [Misc Info](#misc-info)
  - [Command-line flags](#command-line-flags)
  - [Default Ports](#default-ports)
  - [Accessing the Application](#accessing-the-application)
  - [Initial Admin User](#initial-admin-user)
- [Testing](#testing)
  - [Backend](#backend-1)
  - [Frontend](#frontend-1)
  - [Database](#database)
- [Troubleshooting](#troubleshooting)
  - [Frontend build fails](#frontend-build-fails)
  - [Backend build](#backend-build)
  - [Go module issues](#go-module-issues)
  - [Common Runtime Issues](#common-runtime-issues)
    - ["Auth Methods: [password]" instead of "Auth Methods: [chainfs]"](#auth-methods-password-instead-of-chainfs)
    - [Invalid redirect URI](#invalid-redirect-uri)
    - ["User does not exist"](#user-does-not-exist)
    - [Port already in use](#port-already-in-use)
    - [Database locked](#database-locked)
  - [Viewing Logs](#viewing-logs)
- [Building for Production](#building-for-production)
- [Useful Commands](#useful-commands)
- [Additional Resources](#additional-resources)
- [Support](#support)

---

## Prerequisites

**Backend (Go):**
- Go 1.25.0 or higher
- git

**Frontend (Vue.js):**
- Node.js 18+ (LTS recommended)
- npm 9+ or yarn

**Operating System:** 
- Linux

---

## Quick Start (DEV Build)

### Frontend

```bash
cd /home/mem/git/filebrowser/frontend
npm install
npm run build
```

### Backend

note: requires frontend to be built.

```bash
cd /home/mem/git/filebrowser/backend

# Development build (use this by default)
go build -o filebrowser

# Production build (optimized), only use for deployment to azure (not ready yet)
go build -ldflags="-s -w" -o filebrowser

```

### Run 

note: requires frontend and backend to have been built.

```bash
cd /home/mem/git/filebrowser/backend
./filebrowser -c config.dev.yaml
```

Open browser to: `http://localhost:8080`


---

## Running the Application

The application uses YAML configuration files:

- `config.yaml` - Default configuration (uses DEV settings)
- `config.dev.yaml` - DEV environment (points to nansendev.azurewebsites.net)
- `config.uat.yaml` - UAT environment (points to nansenuat.azurewebsites.net)
- `config.prod.yaml` - PROD environment (points to nansenprod.azurewebsites.net)

**Using specific config:**

```bash
cd /home/mem/git/filebrowser/backend

# DEV environment (for all testing use DEV config)
./filebrowser -c config.dev.yaml

# UAT environment
./filebrowser -c config.uat.yaml

# PROD environment
./filebrowser -c config.prod.yaml
```

---

## Misc Info

### Command-line flags

```bash
./filebrowser -h              # Show help
./filebrowser -c              # Print default config
./filebrowser version         # Show version info
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

### Backend

```bash
cd /home/mem/git/filebrowser/backend

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

### Frontend

```bash
cd /home/mem/git/filebrowser/frontend

# Run unit tests
npm run test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch
```

### Database

note: this is used to see if users exist

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

### Frontend build fails

```bash
# Clear node_modules and reinstall
cd /home/mem/git/filebrowser/frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

### Backend build

**fails with "template: pattern matches no files"**

```bash
# Frontend not built yet
cd /home/mem/git/filebrowser/frontend
npm run build

# Then rebuild backend
cd /home/mem/git/filebrowser/backend
go build -o filebrowser
```

### Go module issues

```bash
cd /home/mem/git/filebrowser/backend
go mod tidy
go mod download
go build -o filebrowser
```

### Common Runtime Issues

#### "Auth Methods: [password]" instead of "Auth Methods: [chainfs]"

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

#### Invalid redirect URI" 

from Azure AD B2C

**Problem:** Redirect URI not registered in Azure AD B2C

**Solution:**
- Register callback URL in Azure AD B2C application settings:
  ```
  http://localhost:8080/api/auth/chainfs/callback
  ```
- Ensure exact match (protocol, domain, port, path)

#### "User does not exist"

**Problem:** `createUser: false` but user not in database

**Solution:**
```yaml
# config.dev.yaml
auth:
  methods:
    chainfs:
      createUser: true  # Enable auto-user creation
```

#### Port already in use

**Problem:** Another process using port 8080

**Solution:**
```bash
# its probably just an old copy of filebrowser running

# first, just run killall
killall filebrowser

# check if port is still in use
netstat -tupln | grep 8080

# if a different process ask human to assist
```

#### Database locked

**Problem:** Multiple FileBrowser instances accessing same database

**Solution:**

```bash
# kill all copies of filebrowser
killall filebrowser

# delete files
cd /home/mem/git/filebrowser/backend
rm database.db-shm database.db-wal

# Restart FileBrowser
```

---

## Viewing Logs

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

```bash
# 1. Build frontend (production mode)
cd /home/mem/git/filebrowser/frontend
npm install --production
npm run build

# 2. Build backend (optimized)
cd /home/mem/git/filebrowser/backend
go build -ldflags="-s -w" -o filebrowser

# 3. Prepare distribution
mkdir -p dist
cp filebrowser dist/
cp config.prod.yaml dist/config.yaml

# 4. Create tarball
cd dist
tar -czf filebrowser-chainfs-v1.0.0.tar.gz *
```

---

## Useful Commands

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

For upstream FileBrowser issues:
- Original repository: https://github.com/gtsteffaniak/filebrowser
- Documentation: https://github.com/gtsteffaniak/filebrowser/wiki


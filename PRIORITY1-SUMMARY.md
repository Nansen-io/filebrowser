# Priority 1: ChainFS Authentication - Quick Reference

**Status:** ✅ COMPLETE
**Date:** January 13, 2026
**Commits:** `fb26f204` → `b26ac1b8`

---

## What Was Built

✅ **ChainFS Azure AD B2C OAuth2 authentication**
- Authorization Code Flow (not implicit flow)
- User auto-provisioning
- Role-based admin assignment
- Encrypted token storage (AES-256-GCM)
- Multi-environment support (DEV/UAT/PROD)

---

## Files to Review

### Documentation
- **`CLAUDE.md`** - Complete fork documentation (750 lines)
- **`BUILD.md`** - Build and development guide (691 lines)
- **`IMPLEMENTATION.md`** - This implementation chronicle (you're here!)

### Key Implementation Files
- **`backend/http/chainfs.go`** - Main authentication logic (507 lines)
- **`backend/chainfs/client.go`** - ChainFS API client (60 lines)
- **`backend/config.dev.yaml`** - DEV environment config
- **`frontend/src/views/Login.vue`** - ChainFS login button

### Patch File
- **`PRIORITY1-IMPLEMENTATION.patch`** - Full diff (3,191 lines)

---

## Quick Start

### Build and Run

```bash
# 1. Build frontend (Windows)
cd frontend
npm install
npm run build:windows

# 2. Build backend
cd ../backend
go build -o filebrowser.exe

# 3. Configure (add real client secret)
# Edit config.dev.yaml and add Azure AD B2C client secret

# 4. Run
./filebrowser.exe -c config.dev.yaml
```

### Test Authentication

1. Open `http://localhost:8080`
2. Click "ChainFS Login"
3. Authenticate with Azure AD B2C
4. Should redirect to `/files/` as authenticated user

---

## Configuration Required

### config.dev.yaml
```yaml
auth:
  key: "fb-chainfs-dev-key-32bytes-xxyyz" # Exactly 32 bytes
  methods:
    chainfs:
      enabled: true
      apiBaseUrl: "https://nansendev.azurewebsites.net"
      clientSecret: "YOUR_AZURE_CLIENT_SECRET_HERE" # ← ADD THIS
      createUser: true
      adminClaim: "roles"
      adminClaimValue: "admin"
```

### Azure AD B2C Configuration

**Application ID:** `ae8e4cce-f313-459b-b86b-2fa59b4f1cb8`

**Required Redirect URIs:**
- `http://localhost:8080/api/auth/chainfs/callback` (DEV)
- `https://filebrowser-dev.azurewebsites.net/api/auth/chainfs/callback`
- `https://filebrowser-uat.azurewebsites.net/api/auth/chainfs/callback`
- `https://filebrowser-prod.azurewebsites.net/api/auth/chainfs/callback`

**Client Secret:** Generate in Azure Portal → Certificates & secrets

**Token Claims Required:**
- `sub` - User identifier
- `preferred_username` or `email` - Username
- `roles` or `groups` - For admin assignment

---

## Issues Resolved

| # | Issue | Solution |
|---|-------|----------|
| 1 | Compilation errors | Used `settings.ApplyUserDefaults()` pattern |
| 2 | Auth methods showing [password] | Registered ChainFS in `setupAuth()` |
| 3 | Windows build failing | Created `build-windows.ps1` script |
| 4 | redirect_uri_mismatch | Added redirect URI to Azure AD B2C |
| 5 | Code in URL fragment (#) | Added `response_mode=query` parameter |
| 6 | Client secret required | Added `ClientSecret` to config |
| 7 | Invalid encryption key length | Added 32-byte `auth.key` |
| 8 | User ID = 0 after creation | Reload user after `CreateUser()` |

**Total Debug Time:** ~3 hours
**Total Implementation Time:** ~6.5 hours

---

## Statistics

```
28 files changed
2,275 insertions(+)
112 deletions(-)
3,191 lines in patch file
1,441 lines of documentation
```

---

## Architecture

### Authentication Flow

```
User clicks "ChainFS Login"
  ↓
Backend fetches Azure login URL from ChainFS API
  ↓
Modify redirect_uri to point to FileBrowser callback
  ↓
Redirect user to Azure AD B2C
  ↓
User authenticates
  ↓
Azure redirects to FileBrowser callback with authorization code
  ↓
Backend exchanges code for tokens (access, refresh, ID)
  ↓
Extract user info from ID token
  ↓
Create/update user in database
  - Encrypt Azure tokens (AES-256-GCM)
  - Set LoginMethod = "chainfs"
  ↓
Generate FileBrowser JWT token
  ↓
Set HTTP-only cookie
  ↓
Redirect to /files/
```

### Token Strategy

**Azure AD B2C Tokens** (Server-side, encrypted in database)
- Access token - For ChainFS API calls (Priority 2)
- Refresh token - Auto-refresh when expired
- ID token - User info extraction

**FileBrowser JWT Token** (HTTP-only cookie)
- Session management
- Existing middleware works unchanged
- Expires based on `tokenExpirationHours`

---

## Next Steps

### Priority 2: Right-Click Menu & ChainFS File Sync

**Start Here:**
1. Review `CLAUDE.md` Priority 2 section (lines 50-150)
2. Review ChainFS API endpoints in Swagger
3. Plan database schema for sync metadata
4. Implement right-click menu in frontend
5. Implement ChainFS API client for file operations

**Key Files to Modify:**
- `backend/database/files/` - Add sync metadata fields
- `frontend/src/components/files/` - Add context menu options
- `backend/http/files.go` - Add sync endpoints
- `backend/chainfs/` - Add file operation functions

### Priority 3: Azure Hosting

**Requirements:**
- Docker containerization
- Three instances: DEV, UAT, PROD
- Environment-specific secrets
- Azure Key Vault integration

---

## Useful Commands

### Git
```bash
# View full diff
git diff fb26f204..b26ac1b8

# Apply patch
git apply PRIORITY1-IMPLEMENTATION.patch

# View commit history
git log --oneline fb26f204..b26ac1b8
```

### Build
```bash
# Frontend (Windows)
npm run build:windows

# Frontend (Linux/Mac)
npm run build

# Backend
go build -o filebrowser.exe

# Clean build
go clean && go build -o filebrowser.exe
```

### Testing
```bash
# Run with DEV config
./filebrowser.exe -c config.dev.yaml

# View database users
sqlite3 database.db "SELECT id, username, loginMethod FROM users;"

# Check logs
./filebrowser.exe -c config.dev.yaml 2>&1 | tee filebrowser.log
```

---

## Security Notes

### DO NOT Commit These Files
- `config.dev.yaml` (contains client secret)
- `config.uat.yaml` (contains client secret)
- `config.prod.yaml` (contains client secret and auth key)
- `database.db` (contains user data and encrypted tokens)

### Production Checklist
- [ ] Generate strong random 32-byte `auth.key`
- [ ] Generate new Azure AD B2C client secret
- [ ] Store secrets in Azure Key Vault
- [ ] Use environment variables for sensitive data
- [ ] Enable HTTPS (required for Azure AD B2C)
- [ ] Set `createUser: false` in PROD config
- [ ] Configure Azure Application Insights for monitoring
- [ ] Set up automated backups for database

---

## Support

**Documentation:**
- `CLAUDE.md` - Architecture and API reference
- `BUILD.md` - Build instructions and troubleshooting
- `IMPLEMENTATION.md` - Complete implementation details

**ChainFS API:**
- DEV: https://nansendev.azurewebsites.net/swagger/v1/swagger.json
- UAT: https://nansenuat.azurewebsites.net/swagger/v1/swagger.json
- PROD: https://nansenprod.azurewebsites.net/swagger/v1/swagger.json

**Source Code:**
- FileBrowser: https://github.com/gtsteffaniak/filebrowser
- ChainFS API: `C:\git\azure-blockchain-workbench-app\NasenAPI`

---

**Status:** Ready for Priority 2 implementation 🚀

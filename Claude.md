# FileBrowser Quantum - ChainFS Integration Fork

## Fork Overview

This is a fork of [FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) integrated with the **ChainFS API** - a blockchain-based file storage system with Azure AD B2C authentication. This fork transforms FileBrowser into a hybrid file manager that works with both local files and ChainFS distributed storage.

**Original Repository:** https://github.com/gtsteffaniak/filebrowser

**ChainFS API Environments:**
- **DEV:** https://nansendev.azurewebsites.net/swagger/v1/swagger.json
- **UAT:** https://nansenuat.azurewebsites.net/swagger/v1/swagger.json
- **PROD:** https://nansenprod.azurewebsites.net/swagger/v1/swagger.json

**ChainFS Source Code:** `C:\git\azure-blockchain-workbench-app\NasenAPI`

---

## Why This Fork Exists

The ChainFS project needed a modern, user-friendly web interface for:
1. **Secure Authentication** - Azure AD B2C (not Entra) for enterprise identity management
2. **File Management** - Browsing local files with optional blockchain storage
3. **Audit Trail** - Immutable file versioning through ChainFS blockchain
4. **Multi-Environment Support** - Separate DEV/UAT/PROD deployments

FileBrowser Quantum provides an excellent foundation with its clean architecture, multi-source support, and modern Vue.js frontend.

---

## Integration Goals - Three Priorities

### Priority 1: Authentication ✅ (Current Implementation)

**Goal:** Replace FileBrowser's authentication with ChainFS Azure AD B2C.

**Implementation:**
- Remove password, OIDC, proxy, and no-auth methods
- Integrate Azure AD B2C OAuth2 flow via ChainFS API endpoints
- Auto-provision users on first login
- Support role-based admin assignment via Azure claims
- Dual token system: Azure tokens (for ChainFS API) + FileBrowser JWT (for session management)

**Status:** Implementation in progress

---

### Priority 2: Right-Click Menu & ChainFS File Sync

**Goal:** Allow users to selectively sync files to ChainFS blockchain storage.

**Features:**
- **Right-click menu options:**
  - "Store on ChainFS" - Upload file to blockchain
  - "Update on ChainFS" - Sync changes to existing blockchain file
  - Visual indicators showing sync status

**ChainFS API Endpoints:**
```
/api/NansenFile/DirCreate        - Create directory on ChainFS
/api/NansenFile/DirGetInfo       - Get directory information
/api/NansenFile/DirRename        - Rename directory
/api/NansenFile/DirSubDirs       - List subdirectories
/api/NansenFile/FileCreate       - Create new file on ChainFS
/api/NansenFile/FileExists       - Check if file exists
/api/NansenFile/FileGetDetails   - Get file metadata
/api/NansenFile/FileNewest       - Get latest revision GUID
/api/NansenFile/FileUpdate       - Update existing file
/api/NansenFile/GetFileSimpleInfo - Get file info including genesisGuid
/api/NansenFile/ListOfDirectories - List all directories
/api/NansenFile/SetTags          - Set file tags
/api/Debug/FileEncode            - Reference implementation for encoding files
```

**File Versioning Concepts:**
- **genesisGuid**: The first revision GUID in a file's history (immutable identifier)
- **Current GUID**: The latest revision GUID for a file
- Files on ChainFS maintain complete version history
- Use `/api/NansenFile/GetFileSimpleInfo` to find genesisGuid
- Use `/api/NansenFile/FileNewest` to get latest revision GUID

**Sync Tracking:**
- Store metadata in FileBrowser's BoltDB database:
  - `ChainFsSyncStatus`: "not_synced", "synced", "outdated", "syncing"
  - `ChainFsGenesisGuid`: First revision identifier
  - `ChainFsLatestGuid`: Current revision GUID
  - `ChainFsLastSyncTime`: Timestamp of last sync
- Compare local file modification time vs ChainFS sync time to detect changes
- Visual indicators in UI show sync status per file

**MVP Approach:**
- Start with DEV environment only
- One-way sync: FileBrowser → ChainFS (no download from ChainFS initially)
- Manual sync via right-click menu (not automatic)
- Use FileBrowser file path as ChainFS identifier (or UUID if FileBrowser uses them)

**Status:** Planned (after Priority 1 completion)

---

### Priority 3: Azure Hosting

**Goal:** Deploy three instances of the fork to Azure.

**Deployment Strategy:**
- **Technology:** Docker containers on Azure Container Instances or App Service
- **Three Instances:**
  - DEV → Points to `https://nansendev.azurewebsites.net`
  - UAT → Points to `https://nansenuat.azurewebsites.net`
  - PROD → Points to `https://nansenprod.azurewebsites.net`

**Environment Configuration:**
- Use environment variables or config files to specify ChainFS API URL
- Separate Azure AD B2C redirect URIs per environment
- Environment-specific secrets in Azure Key Vault

**Infrastructure:**
```
Azure Container Registry
  ↓
Azure Container Instance (or App Service)
  ├─ filebrowser-chainfs-dev
  ├─ filebrowser-chainfs-uat
  └─ filebrowser-chainfs-prod
```

**Status:** Planned (after Priority 2 completion)

---

## Architecture Decisions

### Authentication Architecture

**Decision:** Dual Token System

**Rationale:**
1. **Azure AD B2C tokens** are large, short-lived, and specific to ChainFS API
2. **FileBrowser JWT tokens** are lightweight, customizable, and work with existing middleware
3. Separating concerns allows independent token lifecycles
4. Existing FileBrowser session management continues to work unchanged

**Token Flow:**
```
User Login
  ↓
Azure AD B2C issues: ID Token + Access Token + Refresh Token
  ↓
FileBrowser:
  - Stores Azure tokens (encrypted) in database
  - Generates FileBrowser JWT for session management
  - Sets JWT as HTTP-only cookie
  ↓
Subsequent Requests:
  - Frontend: Sends JWT cookie (automatic)
  - FileBrowser: Validates JWT (existing middleware)
  - ChainFS API calls: Uses stored Azure token (decrypted)
```

---

### Storage Architecture

**Decision:** Hybrid Local + Optional ChainFS Sync

**Rationale:**
1. FileBrowser continues managing local files (existing functionality)
2. Users opt-in to blockchain storage per file (not all files need immutability)
3. Local files remain accessible even if ChainFS is unavailable
4. Gradual migration path for existing FileBrowser users

**Storage Model:**
```
Local Filesystem (Primary)
  ├─ File operations work as normal
  ├─ Fast access, familiar UX
  └─ Optional sync markers in database

ChainFS Blockchain (Secondary)
  ├─ Immutable audit trail
  ├─ Version history
  └─ Distributed storage
```

---

### User Management

**Decision:** Auto-Provision Users on First Login

**Configuration Option:** `createUser: true`

**Rationale:**
1. Simplifies onboarding (Azure AD B2C already manages identity)
2. Users don't need separate FileBrowser account creation
3. Admin status determined by Azure claims (centralized role management)
4. Production can disable auto-creation for tighter control

**User Mapping:**
```
Azure AD B2C User
  ↓
Extract: preferred_username, email, roles/groups
  ↓
Create/Update FileBrowser User:
  - Username: From Azure claim
  - LoginMethod: "chainfs"
  - Admin: Based on role claim match
  - AzureAccessToken: Encrypted
  - AzureRefreshToken: Encrypted
```

---

### Token Security

**Decision:** Server-Side Encrypted Storage

**Implementation:**
- Azure tokens never sent to frontend
- Encrypted using existing `config.Auth.Key` before database storage
- Algorithm: AES-256-GCM (or existing FileBrowser encryption method)
- Decrypted only when needed for ChainFS API calls

**Security Benefits:**
1. Tokens not accessible via browser developer tools
2. XSS attacks cannot steal Azure tokens
3. Database compromise still requires decryption key
4. Follows OAuth2 best practices for token storage

---

## Priority 1 Implementation Details

### Authentication Flow (Step-by-Step)

**1. User Initiates Login:**
```
User clicks "ChainFS Login" button
  ↓
Frontend: GET /api/auth/chainfs/login?redirect=/files/documents
```

**2. Backend Fetches Azure Login URL:**
```go
// Backend calls ChainFS API
GET https://nansendev.azurewebsites.net/api/NansenFile/LoginURL

// ChainFS returns:
https://NansenDEV2.b2clogin.com/NansenDEV2.onmicrosoft.com/B2C_1_signupsignin1/oauth2/v2.0/authorize?client_id=ae8e4cce-f313-459b-b86b-2fa59b4f1cb8&redirect_uri=https://jwt.ms/&response_type=token&scope=openid+profile+offline_access...

// Backend modifies redirect_uri:
redirect_uri=https://your-filebrowser.com/api/auth/chainfs/callback

// Backend adds state parameter:
state={nonce}:/files/documents

// Backend redirects user to modified URL
```

**3. User Authenticates with Azure:**
```
Azure AD B2C displays login page
User enters credentials
Multi-factor authentication (if configured)
Azure validates user
```

**4. Azure Redirects to FileBrowser:**
```
GET /api/auth/chainfs/callback?code=ABC123&state={nonce}:/files/documents

Backend:
  - Validates state nonce (CSRF protection)
  - Extracts authorization code
```

**5. Backend Exchanges Code for Tokens:**
```go
// OAuth2 token exchange
POST https://NansenDEV2.b2clogin.com/.../token
  grant_type=authorization_code
  code=ABC123
  redirect_uri=https://your-filebrowser.com/api/auth/chainfs/callback
  client_id=...
  client_secret=...

// Azure responds with:
{
  "access_token": "eyJ0eXAi...",
  "id_token": "eyJ0eXAi...",
  "refresh_token": "...",
  "expires_in": 3600
}
```

**6. Backend Processes User:**
```go
// Parse ID token (JWT)
claims := {
  "sub": "12345678-1234-1234-1234-123456789012",
  "preferred_username": "john.doe@example.com",
  "email": "john.doe@example.com",
  "roles": ["user", "admin"]
}

// Check if user exists
user := database.GetUserByUsername("john.doe@example.com")

if user == nil && config.CreateUser {
  // Create new user
  user = &User{
    Username: "john.doe@example.com",
    LoginMethod: "chainfs",
    Admin: contains(claims["roles"], "admin"),
    AzureAccessToken: encrypt(access_token),
    AzureRefreshToken: encrypt(refresh_token),
    AzureTokenExpiry: now() + expires_in,
  }
  database.SaveUser(user)
} else {
  // Update existing user
  user.AzureAccessToken = encrypt(access_token)
  user.AzureRefreshToken = encrypt(refresh_token)
  user.AzureTokenExpiry = now() + expires_in
  database.UpdateUser(user)
}
```

**7. Backend Generates FileBrowser Session:**
```go
// Generate JWT token (existing FileBrowser logic)
jwtToken := jwt.Sign({
  user: user.ID,
  exp: now() + 2hours,
}, config.Auth.Key)

// Set HTTP-only cookie
setCookie("filebrowser_quantum_jwt", jwtToken, httpOnly: true, secure: true)

// Redirect to original destination
redirect("/files/documents")
```

**8. Subsequent Requests:**
```
All FileBrowser requests include JWT cookie
Existing middleware validates JWT
User object loaded from database
For ChainFS API calls (Priority 2):
  - Decrypt user.AzureAccessToken
  - Check if expired
  - If expired, refresh using user.AzureRefreshToken
  - Make ChainFS API call with Bearer token
```

---

### Code Structure

**Backend Packages:**
```
backend/
├── chainfs/
│   └── client.go              # ChainFS API client functions
├── common/settings/
│   ├── auth.go                # + ChainFsConfig struct
│   └── structs.go             # + ChainFsAuth field
├── database/users/
│   └── users.go               # + LoginMethodChainFs, Azure token fields
└── http/
    ├── chainfs.go             # ChainFS auth handlers (NEW)
    ├── auth.go                # Updated logout handler
    ├── httpRouter.go          # + ChainFS routes
    └── static.go              # + chainfsAvailable flag
```

**Frontend Components:**
```
frontend/
└── src/
    ├── utils/
    │   └── constants.js       # + chainfsAvailable
    └── views/
        └── Login.vue          # + ChainFS login button
```

**Configuration:**
```
backend/
├── config.yaml                # Main config (DEV default)
├── config.dev.yaml            # DEV environment
├── config.uat.yaml            # UAT environment
└── config.prod.yaml           # PROD environment
```

---

### Configuration Schema

**config.yaml:**
```yaml
server:
  port: 8080
  baseURL: "/"

auth:
  tokenExpirationHours: 2
  key: "random-secret-key-for-jwt-signing"
  adminUsername: "admin"
  adminPassword: "admin"

  methods:
    # ChainFS Authentication (NEW)
    chainfs:
      enabled: true
      apiBaseUrl: "https://nansendev.azurewebsites.net"
      createUser: true
      adminClaim: "roles"         # Azure claim containing roles
      adminClaimValue: "admin"    # Value that grants admin

    # Disable other auth methods
    password:
      enabled: false
    oidc:
      enabled: false
    proxy:
      enabled: false
    noauth: false

userDefaults:
  permissions:
    admin: false
    modify: true
    share: true
  locale: "en"
  viewMode: "list"
```

---

## API Reference

### ChainFS Authentication Endpoints

**Get Login URL:**
```bash
curl -X GET 'https://nansendev.azurewebsites.net/api/NansenFile/LoginURL' \
  -H 'accept: text/plain'

# Response:
# https://NansenDEV2.b2clogin.com/NansenDEV2.onmicrosoft.com/B2C_1_signupsignin1/oauth2/v2.0/authorize?client_id=ae8e4cce-f313-459b-b86b-2fa59b4f1cb8&redirect_uri=https://jwt.ms/&response_type=token&scope=openid+profile+offline_access+https://NansenDEV2.onmicrosoft.com/tasks-api/access_as_user&response_mode=fragment
```

**Get Logout URL:**
```bash
curl -X GET 'https://nansendev.azurewebsites.net/api/NansenFile/LogoutURL' \
  -H 'accept: text/plain'

# Response:
# https://nansendev.azurewebsites.net/api/User/Logout
```

### FileBrowser ChainFS Endpoints (NEW)

**Login:**
```
GET /api/auth/chainfs/login?redirect={path}

Response: 302 Redirect to Azure AD B2C
```

**Callback:**
```
GET /api/auth/chainfs/callback?code={code}&state={state}

Response: 302 Redirect to FileBrowser UI (with JWT cookie set)
```

**Logout:**
```
POST /api/auth/chainfs/logout

Response: 200 OK
{
  "logoutRedirectUrl": "https://nansendev.azurewebsites.net/api/User/Logout"
}
```

---

## Development Workflow

### Local Development Setup

**1. Clone the fork:**
```bash
git clone https://github.com/your-username/filebrowser-chainfs.git
cd filebrowser-chainfs
```

**2. Configure for DEV environment:**
```bash
cp backend/config.dev.yaml backend/config.yaml
# Edit config.yaml if needed
```

**3. Start backend:**
```bash
cd backend
go run main.go
# Backend runs on http://localhost:8080
```

**4. Start frontend (development mode):**
```bash
cd frontend
npm install
npm run dev
# Frontend runs on http://localhost:5173 (proxies to backend)
```

**5. Build for production:**
```bash
# Frontend
cd frontend
npm run build

# Backend (includes embedded frontend)
cd backend
go build -o filebrowser
./filebrowser
```

---

### Testing Authentication Flow

**1. Start FileBrowser locally:**
```bash
./filebrowser --config config.dev.yaml
```

**2. Open browser:**
```
http://localhost:8080
```

**3. Click "ChainFS Login"**

**4. Expected flow:**
- Redirects to Azure AD B2C DEV login
- Enter Azure credentials
- Redirects back to FileBrowser
- Logged in, JWT cookie set
- Can browse files

**5. Check database:**
```bash
# User should exist with LoginMethod=chainfs
# Azure tokens should be encrypted in database
```

---

### Environment Switching

**Development:**
```bash
./filebrowser --config config.dev.yaml
# Points to https://nansendev.azurewebsites.net
```

**UAT:**
```bash
./filebrowser --config config.uat.yaml
# Points to https://nansenuat.azurewebsites.net
```

**Production:**
```bash
./filebrowser --config config.prod.yaml
# Points to https://nansenprod.azurewebsites.net
```

---

## Azure AD B2C Configuration Requirements

Before deploying, ensure Azure AD B2C is configured:

**1. Application Registration:**
- Application (client) ID: `ae8e4cce-f313-459b-b86b-2fa59b4f1cb8` (from ChainFS)
- Application type: Web
- User flow: `B2C_1_signupsignin1`

**2. Redirect URIs (must register all):**
```
http://localhost:8080/api/auth/chainfs/callback    (local dev)
https://filebrowser-dev.azurewebsites.net/api/auth/chainfs/callback
https://filebrowser-uat.azurewebsites.net/api/auth/chainfs/callback
https://filebrowser-prod.azurewebsites.net/api/auth/chainfs/callback
```

**3. Token Configuration:**
- Response type: Authorization code flow
- Include ID token: Yes
- Include access token: Yes
- Include refresh token: Yes

**4. Token Claims (must be included in ID token):**
- `sub` - User identifier
- `preferred_username` - Username
- `email` - Email address
- `roles` or `groups` - For admin assignment

**5. API Permissions:**
- Ensure FileBrowser can call ChainFS API with issued tokens
- Scope: `https://NansenDEV2.onmicrosoft.com/tasks-api/access_as_user`

---

## Security Considerations

### HTTPS Requirements
- Azure AD B2C requires HTTPS for redirect URIs (except localhost)
- Production deployment must use HTTPS
- Configure TLS certificates in Azure App Service

### Token Storage
- Azure tokens encrypted using `config.Auth.Key`
- Store key in Azure Key Vault for production
- Never commit keys to git

### CSRF Protection
- State parameter includes nonce
- Validated on callback
- Prevents replay attacks

### Cookie Security
```go
// FileBrowser JWT cookie settings
setCookie("filebrowser_quantum_jwt", token,
  HttpOnly: true,      // No JavaScript access
  Secure: true,        // HTTPS only
  SameSite: "Lax",     // CSRF protection
)
```

### Secrets Management
**Development:**
- Secrets in `config.yaml` (gitignored)

**Production:**
- Use Azure Key Vault
- Environment variables for sensitive data
- Never hardcode secrets

---

## Future Enhancements

### Priority 2 Enhancements
- [ ] Bi-directional sync (download from ChainFS)
- [ ] Automatic sync on file save
- [ ] Conflict resolution UI
- [ ] Bulk sync operations
- [ ] Sync queue with retry logic
- [ ] File diff view (local vs ChainFS)

### Priority 3 Enhancements
- [ ] CI/CD pipeline (GitHub Actions → Azure)
- [ ] Terraform infrastructure as code
- [ ] Application Insights monitoring
- [ ] Azure CDN for static assets
- [ ] Health check endpoints
- [ ] Automated backups

### General Enhancements
- [ ] Audit log of ChainFS operations
- [ ] File provenance viewer (blockchain history)
- [ ] Collaborative file annotations
- [ ] Smart contracts integration
- [ ] Mobile app support
- [ ] SSO with other Azure services

---

## Troubleshooting

### Login Redirect Fails
**Symptom:** "Invalid redirect URI" error from Azure

**Solution:**
1. Check redirect URI is registered in Azure AD B2C
2. Ensure exact match (including protocol, domain, path)
3. Verify `config.Auth.Methods.ChainFsAuth.ApiBaseUrl` is correct

### User Not Created
**Symptom:** Login succeeds but user not in database

**Solution:**
1. Check `createUser: true` in config
2. Verify Azure ID token includes required claims
3. Check backend logs for errors

### Token Expired
**Symptom:** ChainFS API calls fail with 401

**Solution:**
1. Implement token refresh logic (Priority 2)
2. Check `AzureTokenExpiry` in database
3. Verify refresh token is valid

### Cannot Access Files After Login
**Symptom:** Redirected to login page after successful login

**Solution:**
1. Check JWT cookie is set (browser dev tools)
2. Verify `config.Auth.Key` is consistent
3. Check middleware isn't rejecting token

---

## Contributing

This is a private fork for ChainFS integration. Changes should:
1. Maintain compatibility with upstream FileBrowser where possible
2. Follow existing code style (Go backend, Vue frontend)
3. Document ChainFS-specific features clearly
4. Test across DEV/UAT/PROD environments

---

## License

Inherits license from upstream FileBrowser Quantum project.

---

## Contact

For questions about this fork or ChainFS integration:
- ChainFS API Source: `C:\git\azure-blockchain-workbench-app\NasenAPI`
- Original FileBrowser: https://github.com/gtsteffaniak/filebrowser

---

**Document Version:** 1.0
**Last Updated:** 2026-01-07
**Status:** Priority 1 (Authentication) - Implementation in Progress

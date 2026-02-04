# FileBrowser ChainFS Fork

## About

This is a fork of [FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) integrated with ChainFS - a blockchain-based file storage system with Azure AD B2C authentication.

**Original Project:** https://github.com/gtsteffaniak/filebrowser

**ChainFS API Environments:**
- DEV: https://nansendev.azurewebsites.net/swagger/v1/swagger.json
- UAT: https://nansenuat.azurewebsites.net/swagger/v1/swagger.json
- PROD: https://nansenprod.azurewebsites.net/swagger/v1/swagger.json

**ChainFS Source Code:** C:\git\azure-blockchain-workbench-app\NasenAPI

---

## Changes from Original Project

### 1. Custom Theme (Completed)

**Visual Changes:**
- Sidebar background: Teal (#3a7d82) instead of gray
- Header background: Teal (#3a7d82) instead of gray
- Sidebar/header text: Light (#f4f8f8) instead of dark
- Body text: Medium gray (#818793) instead of dark gray
- Links and accents: Teal (#448388) instead of blue
- File icon backgrounds: Light teal with hover states

**Files Modified:**
- `frontend/src/components/sidebar/Sidebar.vue`
- `frontend/src/components/sidebar/SidebarActions.vue`
- `frontend/src/views/bars/Default.vue`
- `frontend/src/components/files/Icon.vue`
- `frontend/public/index.html`
- `frontend/src/css/_variables.css`

**Details:** See `THEME_UPDATES_FINAL.md` for complete color specifications and implementation details.

### 2. Authentication Integration (In Progress)

**Goal:** Replace FileBrowser authentication with Azure AD B2C via ChainFS API.

**ChainFS Auth Endpoints:**
- Login URL: `/api/NansenFile/LoginURL`
- Logout URL: `/api/NansenFile/LogoutURL`

**Authentication Flow:**
1. User clicks "ChainFS Login"
2. Backend fetches Azure AD B2C login URL from ChainFS API
3. User authenticates with Azure AD B2C
4. Azure redirects back with authorization code
5. Backend exchanges code for tokens (access, ID, refresh)
6. Backend creates/updates FileBrowser user with encrypted Azure tokens
7. Backend issues FileBrowser JWT for session management

**Dual Token System:**
- **Azure tokens:** Stored encrypted in database, used for ChainFS API calls
- **FileBrowser JWT:** HTTP-only cookie for session management

**Status:** Implementation in progress (see Todo.md for remaining tasks)

---

## Future Priorities

### Priority 2: ChainFS File Sync (Planned)

Add right-click menu options to sync files to ChainFS blockchain storage:
- "Store on ChainFS" - Upload file to blockchain
- "Update on ChainFS" - Sync changes to existing blockchain file

**Key Concepts:**
- **genesisGuid:** First revision GUID in file history (immutable identifier)
- ChainFS maintains complete file version history
- Local sync metadata stored in FileBrowser's BoltDB database

**Status:** Planned (after authentication completion)

### Priority 3: Azure Hosting (Planned)

Deploy three instances to Azure:
- DEV → Points to ChainFS DEV environment
- UAT → Points to ChainFS UAT environment
- PROD → Points to ChainFS PROD environment

**Status:** Planned (after Priority 2 completion)

---

## Documentation

- **BUILD.md** - Build, run, and troubleshooting instructions
- **Todo.md** - Current tasks and future planning
- **THEME_UPDATES_FINAL.md** - Complete theme implementation details
- **README.md** - Original FileBrowser documentation

---

## Development

See **BUILD.md** for instructions on:
- Building frontend and backend
- Running locally with ChainFS DEV environment
- Configuration options
- Troubleshooting

---

**Fork Status:** Active Development
**Current Focus:** Priority 1 (Authentication Integration)
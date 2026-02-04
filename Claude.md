# Claude Project Instructions

## Git Workflow
- **The user will always handle git commits**
- Never use `git add`, `git commit`, or `git push`
- You can use git commands for reading: `git status`, `git diff`, `git log`
- Focus on making changes, user will commit them

## Documentation Philosophy

### Core Documentation Files
1. **Todo.md** - Future planning and current tasks
2. **Fork.md** - Track changes made from original project (what HAS been done)
3. **BUILD.md** - How-to guide for building, running, and troubleshooting
4. **THEME_UPDATES_FINAL.md** - Theme implementation details

### Rules
- **Do not create planning documents** that persist across sessions
- Planning should only exist in Todo.md or within the current session
- Fork.md documents completed work, not future plans
- Keep documentation concise and avoid redundancy

## Project Context

This is a fork of [FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) integrated with ChainFS blockchain file storage and Azure AD B2C authentication.

**ChainFS API Environments:**
- DEV: https://nansendev.azurewebsites.net
- UAT: https://nansenuat.azurewebsites.net
- PROD: https://nansenprod.azurewebsites.net

**ChainFS Source:** C:\git\azure-blockchain-workbench-app\NasenAPI

## Build Commands

```bash
# Frontend build
cd frontend
npm run build

# Backend build
cd backend
go build -o filebrowser.exe

# Run with config
./filebrowser.exe -c config.dev.yaml
```

## Current Priorities

1. **Priority 1 (In Progress):** Azure AD B2C authentication integration
2. **Priority 2 (Planned):** ChainFS file sync via right-click menu
3. **Priority 3 (Planned):** Azure hosting (DEV/UAT/PROD instances)

See **Todo.md** for detailed task list and **Fork.md** for completed changes.

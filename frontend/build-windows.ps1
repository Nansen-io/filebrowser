# Windows PowerShell build script for FileBrowser frontend

Write-Host "Building FileBrowser Frontend (Windows)..." -ForegroundColor Cyan

# Clean previous builds
Write-Host "Cleaning previous builds..." -ForegroundColor Yellow
Remove-Item -Path "..\backend\http\dist\*" -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item -Path "..\backend\http\embed\*" -Recurse -Force -ErrorAction SilentlyContinue

# Build with Vite
Write-Host "Building with Vite..." -ForegroundColor Yellow
npm run vite-build

# Copy to embed directory
Write-Host "Copying to embed directory..." -ForegroundColor Yellow
New-Item -Path "..\backend\http\embed" -ItemType Directory -Force | Out-Null
Copy-Item -Path "..\backend\http\dist\*" -Destination "..\backend\http\embed\" -Recurse -Force

Write-Host "Build complete!" -ForegroundColor Green
Write-Host "Frontend files are in: backend\http\dist and backend\http\embed" -ForegroundColor Green

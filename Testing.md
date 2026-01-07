  To test the implementation:

  1. Build the backend:
  cd backend
  go build -o filebrowser
  2. Run FileBrowser:
  ./filebrowser --config config.yaml
  3. Access: http://localhost:80
    - You should see the "ChainFS Login" button
    - Clicking it will redirect to Azure AD B2C
    - After authentication, you'll be logged into FileBrowser

  Before deploying, ensure:
  - ✅ Azure AD B2C redirect URI is registered: {your-url}/api/auth/chainfs/callback
  - ✅ ID token includes claims: preferred_username, email, roles
  - ✅ Generate a secure 32-byte key for auth.key in production

  Environment switching:
  - DEV: ./filebrowser --config config.dev.yaml
  - UAT: ./filebrowser --config config.uat.yaml
  - PROD: ./filebrowser --config config.prod.yaml

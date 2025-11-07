# GitHub OAuth Demo (Go)

Overview  
This repository is a minimal demo showing how to implement GitHub OAuth in Go. It uses the standard net/http ServeMux patterns (and simple subrouter-like prefix handlers) to organize routes and demonstrates the full OAuth authorization code flow: redirect to GitHub, user authorizes, GitHub redirects back with a code, server exchanges the code for an access token, and the app displays a basic authenticated profile page.

Prerequisites  
- Go 1.18+ installed  
- A GitHub OAuth App (see setup below)  
- Terminal / browser

GitHub OAuth App setup  
1. In GitHub, go to Settings → Developer settings → OAuth Apps → New OAuth App.  
2. Set "Application name" and "Homepage URL" (e.g. http://localhost:8080).  
3. Set "Authorization callback URL" to: http://localhost:8080/auth/callback  
4. After creating, save the Client ID and Client Secret.

Environment variables  
- GITHUB_CLIENT_ID — Client ID from GitHub  
- GITHUB_CLIENT_SECRET — Client Secret from GitHub  
- OAUTH_REDIRECT_URL — (optional) should match callback URL (default: http://localhost:8080/auth/callback)  
- SESSION_KEY — secret for session/cookie signing (or keep simple for demo)  
- PORT — (optional) TCP port, default 8080

Build and run  
1. From the project root (c:\Users\PROTECH_WD\Documents\GO\tst):  
   go build -o bin/server  
   ./bin/server  
   or:  
   go run main.go

2. Open a browser and visit:  
   http://localhost:8080/

Main routes (example)  
- GET  /                 — Home / welcome  
- GET  /auth/login       — Starts OAuth: redirects to GitHub authorize URL  
- GET  /auth/callback    — OAuth callback: exchange code for token, create session  
- GET  /profile          — Protected: shows basic user info (if authenticated)  
- GET  /auth/logout      — Clears session / logs out

Typical flow  
1. User clicks "Login with GitHub" (or visits /auth/login).  
2. Server redirects to GitHub's authorization endpoint with client_id, redirect_uri, scope, state.  
3. User authorizes; GitHub redirects to /auth/callback?code=...&state=...  
4. Server validates state, exchanges code for an access token, optionally fetches user info, then stores auth info in session/cookie.  
5. User can visit /profile to see their username/email fetched from GitHub.

Examples  
- Start login in a browser:  
  Open http://localhost:8080/auth/login

- Check health/home:  
  curl http://localhost:8080/

Security & notes  
- Do not commit Client Secret or SESSION_KEY to source control. Use environment variables.  
- Callback/redirect URL must match the URL configured in your GitHub OAuth app exactly (including protocol and port).  
- For production, use HTTPS for redirect URLs and secure cookies.  
- Consider using a package for sessions (gorilla/sessions) and proper CSRF/state handling for production.  
- For more advanced routing or URL params, a third-party router (gorilla/mux, chi) can be used; this demo keeps dependencies minimal.

Troubleshooting  
- "Invalid redirect_uri" from GitHub: ensure the callback URL registered on GitHub matches the URL your app sends.  
- Missing code param at callback: ensure GitHub redirected correctly and state handling is correct.  
- Port in use: change PORT or stop the process using the port.

License  
MIT (or change as appropriate)
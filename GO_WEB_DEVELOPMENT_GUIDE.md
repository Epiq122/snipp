# Go Web Development Project Guide

This document serves as a comprehensive guide for Go web application development, based on the Snippet project. It
outlines key components, best practices, and implementation patterns that can be reused in future projects.

## Project Structure

A well-organized project structure helps maintain code quality and promotes separation of concerns:

```
cmd/web/               # Application entry point and web components
  ├─ main.go           # App bootstrap, configuration, and dependency injection
  ├─ routes.go         # HTTP route definitions with middleware chains
  ├─ handlers.go       # HTTP request handlers with form processing and sessions
  ├─ middleware.go     # HTTP middleware functions
  ├─ templates.go      # Template functions and cache with flash support
  └─ helpers.go        # Shared helper functions, form decoding, and session helpers
internal/models/       # Data models and database operations
  ├─ snippets.go       # Core model definitions with CRUD operations
  └─ errors.go         # Custom error types
internal/validator/    # Input validation framework
  └─ validator.go      # Reusable validation functions and validator struct
tls/                   # TLS certificate management
  ├─ cert.pem          # TLS certificate (development)
  └─ key.pem           # TLS private key (development)
ui/                    # User interface components
  ├─ html/             # HTML templates
  │   ├─ base.tmpl     # Base layout template with flash message support
  │   ├─ pages/        # Page-specific templates (home, view, create)
  │   └─ partials/     # Reusable template components (nav)
  └─ static/           # Static assets
      ├─ css/          # Stylesheets with form and flash message styling
      ├─ js/           # JavaScript files
      └─ img/          # Images and icons
```

## Core Components and Implementation Patterns

### 1. HTTPS/TLS Server Configuration

#### Complete TLS Server Implementation

```go
// Enhanced application bootstrap with TLS support
func main() {
addr := flag.String("addr", ":8080", "http service address")
dsn := flag.String("dsn", "web:%s@/snippetbox?parseTime=true", "MySQL data source name")
flag.Parse()

logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// Database connection setup
// ... existing database setup code ...

// Session manager with secure cookies for HTTPS
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour
sessionManager.Cookie.Secure = true // HTTPS-only cookies

app := &application{
logger:         logger,
snippets:       &models.SnippetModel{DB: db},
templateCache:  templateCache,
formDecoder:    formDecoder,
sessionManager: sessionManager,
}

// Modern TLS configuration
tlsConfig := &tls.Config{
CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
}

// Professional server configuration
srv := &http.Server{
Addr:         *addr,
Handler:      app.routes(),
ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
TLSConfig:    tlsConfig,
IdleTimeout:  time.Minute,
ReadTimeout:  5 * time.Second,
WriteTimeout: 10 * time.Second,
}

logger.Info("starting on server", "addr", *addr)
err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
logger.Error(err.Error())
os.Exit(1)
}
```

#### TLS Configuration Best Practices

```go
// Production-ready TLS configuration
tlsConfig := &tls.Config{
// Use modern elliptic curves for better performance and security
CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},

// Additional production settings
MinVersion: tls.VersionTLS12, // Minimum TLS 1.2
CipherSuites: []uint16{
tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
},
}
```

#### Server Timeout Configuration

```go
// Professional server configuration with timeouts
srv := &http.Server{
Addr:     *addr,
Handler:  app.routes(),
ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
TLSConfig: tlsConfig,

// Connection timeouts for production reliability
IdleTimeout:  time.Minute,      // How long to keep connections open
ReadTimeout:  5 * time.Second,  // Time to read request headers/body
WriteTimeout: 10 * time.Second, // Time to write response
}
```

### 2. Certificate Management Patterns

#### Development Certificate Structure

```
tls/
├─ cert.pem    # Self-signed certificate for localhost
└─ key.pem     # Private key for the certificate
```

#### Certificate Generation for Development

```bash
# Generate self-signed certificate for development
go run $GOROOT/src/crypto/tls/generate_cert.go --host=localhost
```

#### Production Certificate Deployment

```go
// Environment-based certificate paths for production
certFile := os.Getenv("TLS_CERT_PATH")
keyFile := os.Getenv("TLS_KEY_PATH")

if certFile == "" {
certFile = "./tls/cert.pem" // Development default
}
if keyFile == "" {
keyFile = "./tls/key.pem" // Development default
}

err = srv.ListenAndServeTLS(certFile, keyFile)
```

### 3. Enhanced Security Configuration

#### Secure Session Cookies

```go
// Session configuration for HTTPS
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour

// HTTPS-only security settings
sessionManager.Cookie.Secure = true // HTTPS only
sessionManager.Cookie.HttpOnly = true // No JavaScript access
sessionManager.Cookie.SameSite = http.SameSiteStrictMode // CSRF protection
```

#### Security Headers Integration

```go
// Enhanced security headers for HTTPS
func commonHeaders(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
// Enhanced CSP for HTTPS
w.Header().Set("Content-Security-Policy",
"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

// HTTPS security headers
w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "deny")
w.Header().Set("X-XSS-Protection", "0")
w.Header().Set("Server", "Go")

next.ServeHTTP(w, r)
})
}
```

### 2. Session Management System

#### Session Manager Configuration

```go
// Session manager setup with MySQL backend
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour

// Key features:
// - Database-backed session storage for scalability
// - 12-hour automatic expiration
// - Secure session cookies
// - Integration with existing database connection
```

#### Session Middleware Integration

```go
// Dynamic middleware chain for session-enabled routes
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

// Static files don't need sessions
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// Dynamic routes with session middleware
dynamic := alice.New(app.sessionManager.LoadAndSave)
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

// Standard middleware chain for all routes
standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
return standard.Then(mux)
}
```

### 3. Flash Messaging System

#### Flash Message Implementation

```go
// Template data structure with flash support
type templateData struct {
CurrentYear int
Snippet     models.Snippet
Snippets    []models.Snippet
Form        any
Flash       string  // Added for flash messages
}

// Helper function automatically populates flash messages
func (app *application) newTemplateData(r *http.Request) templateData {
return templateData{
CurrentYear: time.Now().Year(),
Flash:       app.sessionManager.PopString(r.Context(), "flash"),
}
}

// Setting flash messages in handlers
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
// ... validation and processing ...

// Set success flash message
app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

#### Template Integration for Flash Messages

```html
<!-- Base template with flash message display -->
{{define "base"}}
<!doctype html>
<html lang="en">
<head>
    <!-- head content -->
</head>
<body>
    <header>
        <h1><a href="/">Snippetbox</a></h1>
    </header>
    {{template "nav" .}}
    <main>
        {{with .Flash}}
        <div class="flash">{{.}}</div>
        {{end}}
        {{template "main" .}}
    </main>
    <!-- footer and scripts -->
</body>
</html>
{{end}}
```

#### Flash Message Styling

```css
/* Professional flash message styling */
div.flash {
    color: #FFFFFF;
    font-weight: bold;
    background-color: #34495E;
    padding: 18px;
    margin-bottom: 36px;
    text-align: center;
}

/* Error message styling for contrast */
div.error {
    color: #FFFFFF;
    background-color: #C0392B;
    padding: 18px;
    margin-bottom: 36px;
    font-weight: bold;
    text-align: center;
}
```

### 4. Advanced Middleware Architecture

#### Sophisticated Middleware Composition

```go
// Two-tier middleware architecture
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

// Static routes - no sessions needed
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// Dynamic routes - with session support
dynamic := alice.New(app.sessionManager.LoadAndSave)
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

// Standard middleware - applied to all routes
standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

return standard.Then(mux)
}
```

#### Benefits of This Architecture

- **Performance**: Static files bypass session processing
- **Security**: Session handling only where needed
- **Maintainability**: Clear separation of middleware concerns
- **Scalability**: Efficient resource usage

### 5. Application Configuration and Bootstrap

```go
// Command-line flag parsing for configuration
addr := flag.String("addr", ":8080", "HTTP network address")
dsn := flag.String("dsn", "web:%s@/snippetbox?parseTime=true", "MySQL data source name")
flag.Parse()

// Application struct for dependency injection
type application struct {
logger         *slog.Logger
snippets       *models.SnippetModel
templateCache  map[string]*template.Template
formDecoder    *form.Decoder
sessionManager *scs.SessionManager // Added for session management
}

// Environment variable handling for sensitive data
password := os.Getenv("DB_PASSWORD")
if password == "" {
logger.Error("DB_PASSWORD environment variable not set")
os.Exit(1)
}

// Form decoder initialization
formDecoder := form.NewDecoder()

// Session manager initialization with database storage
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour
```

## Core Components and Implementation Patterns

### 1. Application Configuration and Bootstrap

```go
// Command-line flag parsing for configuration
addr := flag.String("addr", ":8080", "HTTP network address")
dsn := flag.String("dsn", "web:%s@/snippetbox?parseTime=true", "MySQL data source name")
flag.Parse()

// Application struct for dependency injection
type application struct {
logger         *slog.Logger
snippets       *models.SnippetModel
templateCache  map[string]*template.Template
formDecoder    *form.Decoder
sessionManager *scs.SessionManager // Added for session management
}

// Environment variable handling for sensitive data
password := os.Getenv("DB_PASSWORD")
if password == "" {
logger.Error("DB_PASSWORD environment variable not set")
os.Exit(1)
}

// Form decoder initialization
formDecoder := form.NewDecoder()

// Session manager initialization with database storage
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour
```

### 2. Session Management System

#### Session Manager Configuration

```go
// Session manager setup with MySQL backend
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour

// Key features:
// - Database-backed session storage for scalability
// - 12-hour automatic expiration
// - Secure session cookies
// - Integration with existing database connection
```

#### Session Middleware Integration

```go
// Dynamic middleware chain for session-enabled routes
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

// Static files don't need sessions
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// Dynamic routes with session middleware
dynamic := alice.New(app.sessionManager.LoadAndSave)
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

// Standard middleware chain for all routes
standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
return standard.Then(mux)
}
```

### 3. Flash Messaging System

#### Flash Message Implementation

```go
// Template data structure with flash support
type templateData struct {
CurrentYear int
Snippet     models.Snippet
Snippets    []models.Snippet
Form        any
Flash       string  // Added for flash messages
}

// Helper function automatically populates flash messages
func (app *application) newTemplateData(r *http.Request) templateData {
return templateData{
CurrentYear: time.Now().Year(),
Flash:       app.sessionManager.PopString(r.Context(), "flash"),
}
}

// Setting flash messages in handlers
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
// ... validation and processing ...

// Set success flash message
app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

#### Template Integration for Flash Messages

```html
<!-- Base template with flash message display -->
{{define "base"}}
<!doctype html>
<html lang="en">
<head>
    <!-- head content -->
</head>
<body>
    <header>
        <h1><a href="/">Snippetbox</a></h1>
    </header>
    {{template "nav" .}}
    <main>
        {{with .Flash}}
        <div class="flash">{{.}}</div>
        {{end}}
        {{template "main" .}}
    </main>
    <!-- footer and scripts -->
</body>
</html>
{{end}}
```

#### Flash Message Styling

```css
/* Professional flash message styling */
div.flash {
    color: #FFFFFF;
    font-weight: bold;
    background-color: #34495E;
    padding: 18px;
    margin-bottom: 36px;
    text-align: center;
}

/* Error message styling for contrast */
div.error {
    color: #FFFFFF;
    background-color: #C0392B;
    padding: 18px;
    margin-bottom: 36px;
    font-weight: bold;
    text-align: center;
}
```

### 4. Advanced Middleware Architecture

#### Sophisticated Middleware Composition

```go
// Two-tier middleware architecture
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

// Static routes - no sessions needed
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// Dynamic routes - with session support
dynamic := alice.New(app.sessionManager.LoadAndSave)
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

// Standard middleware - applied to all routes
standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

return standard.Then(mux)
}
```

#### Benefits of This Architecture

- **Performance**: Static files bypass session processing
- **Security**: Session handling only where needed
- **Maintainability**: Clear separation of middleware concerns
- **Scalability**: Efficient resource usage

### 5. Application Configuration and Bootstrap

```go
// Command-line flag parsing for configuration
addr := flag.String("addr", ":8080", "HTTP network address")
dsn := flag.String("dsn", "web:%s@/snippetbox?parseTime=true", "MySQL data source name")
flag.Parse()

// Application struct for dependency injection
type application struct {
logger         *slog.Logger
snippets       *models.SnippetModel
templateCache  map[string]*template.Template
formDecoder    *form.Decoder
sessionManager *scs.SessionManager // Added for session management
}

// Environment variable handling for sensitive data
password := os.Getenv("DB_PASSWORD")
if password == "" {
logger.Error("DB_PASSWORD environment variable not set")
os.Exit(1)
}

// Form decoder initialization
formDecoder := form.NewDecoder()

// Session manager initialization with database storage
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour
```

## Core Components and Implementation Patterns

### 1. Application Configuration and Bootstrap

```go
// Command-line flag parsing for configuration
addr := flag.String("addr", ":8080", "HTTP network address")
dsn := flag.String("dsn", "web:%s@/snippetbox?parseTime=true", "MySQL data source name")
flag.Parse()

// Application struct for dependency injection
type application struct {
logger         *slog.Logger
snippets       *models.SnippetModel
templateCache  map[string]*template.Template
formDecoder    *form.Decoder
sessionManager *scs.SessionManager // Added for session management
}

// Environment variable handling for sensitive data
password := os.Getenv("DB_PASSWORD")
if password == "" {
logger.Error("DB_PASSWORD environment variable not set")
os.Exit(1)
}

// Form decoder initialization
formDecoder := form.NewDecoder()

// Session manager initialization with database storage
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour
```

### 2. Session Management System

#### Session Manager Configuration

```go
// Session manager setup with MySQL backend
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour

// Key features:
// - Database-backed session storage for scalability
// - 12-hour automatic expiration
// - Secure session cookies
// - Integration with existing database connection
```

#### Session Middleware Integration

```go
// Dynamic middleware chain for session-enabled routes
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

// Static files don't need sessions
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// Dynamic routes with session middleware
dynamic := alice.New(app.sessionManager.LoadAndSave)
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

// Standard middleware chain for all routes
standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
return standard.Then(mux)
}
```

### 3. Flash Messaging System

#### Flash Message Implementation

```go
// Template data structure with flash support
type templateData struct {
CurrentYear int
Snippet     models.Snippet
Snippets    []models.Snippet
Form        any
Flash       string  // Added for flash messages
}

// Helper function automatically populates flash messages
func (app *application) newTemplateData(r *http.Request) templateData {
return templateData{
CurrentYear: time.Now().Year(),
Flash:       app.sessionManager.PopString(r.Context(), "flash"),
}
}

// Setting flash messages in handlers
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
// ... validation and processing ...

// Set success flash message
app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

#### Template Integration for Flash Messages

```html
<!-- Base template with flash message display -->
{{define "base"}}
<!doctype html>
<html lang="en">
<head>
    <!-- head content -->
</head>
<body>
    <header>
        <h1><a href="/">Snippetbox</a></h1>
    </header>
    {{template "nav" .}}
    <main>
        {{with .Flash}}
        <div class="flash">{{.}}</div>
        {{end}}
        {{template "main" .}}
    </main>
    <!-- footer and scripts -->
</body>
</html>
{{end}}
```

#### Flash Message Styling

```css
/* Professional flash message styling */
div.flash {
    color: #FFFFFF;
    font-weight: bold;
    background-color: #34495E;
    padding: 18px;
    margin-bottom: 36px;
    text-align: center;
}

/* Error message styling for contrast */
div.error {
    color: #FFFFFF;
    background-color: #C0392B;
    padding: 18px;
    margin-bottom: 36px;
    font-weight: bold;
    text-align: center;
}
```

### 4. Advanced Middleware Architecture

#### Sophisticated Middleware Composition

```go
// Two-tier middleware architecture
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

// Static routes - no sessions needed
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// Dynamic routes - with session support
dynamic := alice.New(app.sessionManager.LoadAndSave)
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

// Standard middleware - applied to all routes
standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

return standard.Then(mux)
}
```

#### Benefits of This Architecture

- **Performance**: Static files bypass session processing
- **Security**: Session handling only where needed
- **Maintainability**: Clear separation of middleware concerns
- **Scalability**: Efficient resource usage

## HTTPS/TLS Implementation Guide

### 1. Development Environment Setup

#### Self-Signed Certificate Creation

```bash
# Using Go's built-in certificate generator
go run $GOROOT/src/crypto/tls/generate_cert.go --host=localhost

# Or using openssl
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

#### Browser Configuration for Development

When using self-signed certificates in development:

1. **Chrome/Edge**: Navigate to `chrome://flags/#allow-insecure-localhost` and enable
2. **Firefox**: Click "Advanced" → "Accept the Risk and Continue"
3. **curl**: Use `-k` flag to accept self-signed certificates

### 2. Production Deployment Considerations

#### Let's Encrypt Integration

```go
// Production setup with Let's Encrypt
import "golang.org/x/crypto/acme/autocert"

// Automatic certificate management
m := &autocert.Manager{
Cache:      autocert.DirCache("certs"),
Prompt:     autocert.AcceptTOS,
HostPolicy: autocert.HostWhitelist("yourdomain.com"),
}

srv := &http.Server{
Addr:      ":443",
Handler:   app.routes(),
TLSConfig: m.TLSConfig(),
}

go http.ListenAndServe(":80", m.HTTPHandler(nil)) // HTTP->HTTPS redirect
srv.ListenAndServeTLS("", "")
```

#### Load Balancer Considerations

```go
// When behind a load balancer with TLS termination
func (app *application) secureHeaders(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
// Check for forwarded protocol header
if r.Header.Get("X-Forwarded-Proto") == "https" {
w.Header().Set("Strict-Transport-Security", "max-age=31536000")
}
next.ServeHTTP(w, r)
})
}
```

### 3. Security Best Practices

#### TLS Configuration Hardening

```go
// Production-hardened TLS configuration
tlsConfig := &tls.Config{
MinVersion: tls.VersionTLS12,
MaxVersion: tls.VersionTLS13,
CurvePreferences: []tls.CurveID{
tls.X25519,
tls.CurveP256,
},
PreferServerCipherSuites: true,
CipherSuites: []uint16{
tls.TLS_AES_256_GCM_SHA384,
tls.TLS_AES_128_GCM_SHA256,
tls.TLS_CHACHA20_POLY1305_SHA256,
tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
},
}
```

#### Certificate Monitoring and Rotation

```go
// Certificate expiration monitoring
func (app *application) checkCertificateExpiry() {
cert, err := tls.LoadX509KeyPair("./tls/cert.pem", "./tls/key.pem")
if err != nil {
app.logger.Error("failed to load certificate", "error", err)
return
}

x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
if err != nil {
app.logger.Error("failed to parse certificate", "error", err)
return
}

daysUntilExpiry := time.Until(x509Cert.NotAfter).Hours() / 24
if daysUntilExpiry < 30 {
app.logger.Warn("certificate expiring soon",
"days_remaining", int(daysUntilExpiry),
"expires_at", x509Cert.NotAfter,
)
}
}
```

## Performance and Reliability Enhancements

### 1. Connection Management

```go
// Optimized server configuration for production
srv := &http.Server{
Addr:    *addr,
Handler: app.routes(),

// Timeouts prevent resource exhaustion
ReadTimeout:       5 * time.Second,
ReadHeaderTimeout: 2 * time.Second,
WriteTimeout:      10 * time.Second,
IdleTimeout:       time.Minute,

// Connection limits
MaxHeaderBytes: 1 << 20, // 1MB
}
```

### 2. Graceful Shutdown

```go
// Graceful server shutdown implementation
func main() {
// ... server setup ...

// Channel to listen for interrupt signal
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

// Start server in goroutine
go func () {
logger.Info("starting server", "addr", *addr)
if err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil && err != http.ErrServerClosed {
logger.Error("server error", "error", err)
os.Exit(1)
}
}()

// Wait for interrupt signal
<-stop
logger.Info("shutting down server...")

// Create context with timeout for shutdown
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Attempt graceful shutdown
if err := srv.Shutdown(ctx); err != nil {
logger.Error("server forced to shutdown", "error", err)
}

logger.Info("server exited")
}
```

## Security Monitoring and Logging

### 1. TLS Connection Logging

```go
// Enhanced logging for TLS connections
func (app *application) logTLSConnection(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
if r.TLS != nil {
app.logger.Info("TLS connection",
"version", tlsVersionString(r.TLS.Version),
"cipher_suite", tls.CipherSuiteName(r.TLS.CipherSuite),
"server_name", r.TLS.ServerName,
)
}
next.ServeHTTP(w, r)
})
}

func tlsVersionString(version uint16) string {
switch version {
case tls.VersionTLS10:
return "TLS 1.0"
case tls.VersionTLS11:
return "TLS 1.1"
case tls.VersionTLS12:
return "TLS 1.2"
case tls.VersionTLS13:
return "TLS 1.3"
default:
return "Unknown"
}
}
```

### 2. Security Event Monitoring

```go
// Security event logging for HTTPS applications
func (app *application) logSecurityEvent(r *http.Request, event string, details map[string]interface{}) {
logData := map[string]interface{}{
"event":      event,
"ip":         r.RemoteAddr,
"user_agent": r.UserAgent(),
"method":     r.Method,
"uri":        r.URL.RequestURI(),
"timestamp":  time.Now().Unix(),
}

// Add TLS information if available
if r.TLS != nil {
logData["tls_version"] = tlsVersionString(r.TLS.Version)
logData["cipher_suite"] = tls.CipherSuiteName(r.TLS.CipherSuite)
}

// Merge additional details
for k, v := range details {
logData[k] = v
}

app.logger.Info("security event", slog.Any("details", logData))
}
```

## Testing HTTPS Applications

### 1. Testing with Self-Signed Certificates

```go
// Test client configuration for self-signed certificates
func createTestClient() *http.Client {
tr := &http.Transport{
TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
return &http.Client{Transport: tr}
}

// Integration test example
func TestHTTPSEndpoint(t *testing.T) {
client := createTestClient()

resp, err := client.Get("https://localhost:8080/")
if err != nil {
t.Fatalf("Failed to make HTTPS request: %v", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
t.Errorf("Expected status 200, got %d", resp.StatusCode)
}
}
```

### 2. Certificate Validation Testing

```go
// Test certificate loading and validation
func TestCertificateLoading(t *testing.T) {
cert, err := tls.LoadX509KeyPair("./tls/cert.pem", "./tls/key.pem")
if err != nil {
t.Fatalf("Failed to load certificate: %v", err)
}

// Validate certificate can be parsed
x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
if err != nil {
t.Fatalf("Failed to parse certificate: %v", err)
}

// Check certificate hasn't expired
if time.Now().After(x509Cert.NotAfter) {
t.Error("Certificate has expired")
}

// Check certificate is valid for localhost
err = x509Cert.VerifyHostname("localhost")
if err != nil {
t.Errorf("Certificate not valid for localhost: %v", err)
}
}
```

## Deployment Considerations

### 1. Environment Configuration

```go
// Environment-aware certificate management
func getCertificatePaths() (string, string) {
env := os.Getenv("APP_ENV")

switch env {
case "production":
return os.Getenv("TLS_CERT_PATH"), os.Getenv("TLS_KEY_PATH")
case "staging":
return "./certs/staging-cert.pem", "./certs/staging-key.pem"
default: // development
return "./tls/cert.pem", "./tls/key.pem"
}
}
```

### 2. Health Checks for HTTPS

```go
// Health check endpoint for HTTPS applications
func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")

health := map[string]interface{}{
"status":    "healthy",
"timestamp": time.Now().Unix(),
"version":   "0.8.0",
"tls":       r.TLS != nil,
}

if r.TLS != nil {
health["tls_version"] = tlsVersionString(r.TLS.Version)
}

json.NewEncoder(w).Encode(health)
}
```

## Conclusion

The HTTPS/TLS implementation represents a crucial milestone in creating production-ready web applications with Go. The
implementation provides:

- **End-to-End Security**: All HTTP traffic encrypted with TLS
- **Modern Cryptographic Standards**: Support for latest TLS versions and cipher suites
- **Production-Ready Configuration**: Proper timeouts, error handling, and security headers
- **Development-Friendly Setup**: Self-signed certificates for local development
- **Scalable Architecture**: Configuration that works from development to production

This foundation enables deploying secure web applications that meet modern security standards and can handle production
workloads safely.

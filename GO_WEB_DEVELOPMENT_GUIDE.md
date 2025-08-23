# Go Web Development Project Guide

This document serves as a comprehensive guide for Go web application development, based on the Snippet project. It
outlines key components, best practices, and implementation patterns that can be reused in future projects.

## Project Structure

A well-organized project structure helps maintain code quality and promotes separation of concerns:

```
cmd/web/               # Application entry point and web components
  ├─ main.go           # App bootstrap, configuration, and dependency injection
  ├─ routes.go         # HTTP route definitions with middleware chains and authentication
  ├─ handlers.go       # HTTP request handlers with authentication and form processing
  ├─ middleware.go     # HTTP middleware functions including authentication and CSRF
  ├─ templates.go      # Template functions and cache with authentication support
  ├─ helpers.go        # Shared helper functions, form decoding, and authentication helpers
  └─ context.go        # Context keys and authentication state management
internal/models/       # Data models and database operations
  ├─ snippets.go       # Core snippet model with CRUD operations
  ├─ users.go          # User model with authentication operations and existence validation
  └─ errors.go         # Custom error types including authentication errors
internal/validator/    # Input validation framework
  └─ validator.go      # Reusable validation functions with authentication support
tls/                   # TLS certificate management
  ├─ cert.pem          # TLS certificate (development)
  └─ key.pem           # TLS private key (development)
ui/                    # User interface components
  ├─ html/             # HTML templates
  │   ├─ base.tmpl     # Base layout template with authentication state
  │   ├─ pages/        # Page-specific templates (home, view, create, signup, login)
  │   └─ partials/     # Reusable template components (nav with authentication)
  └─ static/           # Static assets
      ├─ css/          # Stylesheets with form and authentication styling
      ├─ js/           # JavaScript files
      └─ img/          # Images and icons
```

## Core Components and Implementation Patterns

### 1. Advanced Context-Based Authentication System

#### Context Key Definition and Management

```go
// context.go - Authentication context management
type contextKey string
const isAuthenticatedContextKey = contextKey("isAuthenticated")
```

#### Enhanced Authentication Middleware

```go
// Advanced authentication middleware with user existence validation
func (app *application) authenticate(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
if id == 0 {
next.ServeHTTP(w, r)
return
}

// Validate user still exists in database
exists, err := app.users.Exists(id)
if err != nil {
app.serverError(w, r, err)
return
}

if exists {
ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
r = r.WithContext(ctx)
}
next.ServeHTTP(w, r)
})
}
```

#### Context-Aware Authentication Helper

```go
// Enhanced isAuthenticated helper using context values
func (app *application) isAuthenticated(r *http.Request) bool {
isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
if !ok {
return false
}
return isAuthenticated
}
```

### 2. User Existence Verification System

#### Efficient Database User Validation

```go
// UserModel with existence validation
func (m *UserModel) Exists(id int) (bool, error) {
var exists bool
stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`
err := m.DB.QueryRow(stmt, id).Scan(&exists)
return exists, err
}
```

### 3. Enhanced Middleware Architecture

#### Multi-Layered Middleware Integration

```go
// Sophisticated middleware composition with context integration
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

// Static files (no middleware needed)
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// Dynamic routes with session, CSRF protection, and authentication
dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF, app.authenticate)
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

// Public authentication routes
mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

// Protected routes requiring authentication
protected := dynamic.Append(app.requireAuthentication)
mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

// Standard middleware for all routes
standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
return standard.Then(mux)
}
```

## Advanced Authentication Patterns

### 1. Context-Based State Management

The authentication system now uses Go's context package for superior state management:

```go
// Benefits of context-based authentication:
// - Thread-safe authentication state
// - Request-scoped authentication data
// - Clean separation from session management
// - Improved testability and maintainability
```

### 2. User Existence Validation

Every authenticated request now validates user existence:

```go
// Protection against stale sessions:
// - Automatic validation on each request
// - Protection against deleted users with active sessions
// - Efficient database queries with SELECT EXISTS()
// - Graceful handling of database errors
```

### 3. Multi-Layered Security

The middleware architecture provides comprehensive security:

```go
// Security layers:
// 1. Session management (LoadAndSave)
// 2. CSRF protection (preventCSRF)
// 3. User authentication (authenticate)
// 4. Route protection (requireAuthentication)
```

## Project Evolution Summary

The Snippet project has now evolved through these major milestones:

- **v0.1.0**: Basic HTTP server
- **v0.2.0**: Template rendering and static assets
- **v0.3.0**: Structured logging and error handling
- **v0.4.0**: MySQL database integration
- **v0.5.0**: Professional middleware system
- **v0.6.0**: Complete form handling and validation
- **v0.7.0**: Session management and flash messaging
- **v0.8.0**: HTTPS/TLS security implementation
- **v0.9.0**: Complete user authentication system
- **v0.10.0**: Advanced context-based authentication *(NEW)*

## Security Enhancements in v0.10.0

### 1. Enhanced Authentication Security

- Context-based authentication state management
- Automatic user existence validation on each request
- Protection against deleted users with active sessions
- Improved separation of concerns between session and authentication

### 2. Performance Improvements

- Context-based authentication state caching
- Efficient database queries for user validation
- Reduced session database calls through context caching

### 3. Code Quality Improvements

- Clean separation of authentication concerns
- Improved testability through context pattern
- Better maintainability with modular middleware design

## Dependencies and Libraries

### Current Dependencies

- `golang.org/x/net/context` - Context package for Go (now part of standard library in newer versions)
- `github.com/alexedwards/scs/v2` - Session management
- `github.com/justinas/nosurf` - CSRF protection
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/justinas/alice` - Middleware chaining
- `github.com/go-playground/form/v4` - Form processing

### Integration Benefits

- Professional authentication with context management
- Secure session handling with user validation
- Comprehensive middleware architecture
- Production-ready security standards

## Conclusion

Version 0.10.0 represents a significant advancement in authentication architecture, moving from basic session-based
authentication to a sophisticated context-aware system. The implementation provides:

- **Advanced Security**: Multi-layered authentication with user existence validation
- **Superior Architecture**: Context-based state management following Go best practices
- **Enhanced Performance**: Efficient authentication state caching and database queries
- **Production-Ready**: Enterprise-level authentication patterns and security measures

This foundation enables building complex, secure web applications that can scale while maintaining the highest security
and performance standards.

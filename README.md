# Snippet

A simple Go web application for creating and viewing text snippets. This repository is being developed incrementally,
with each section documented and versioned so that viewers can follow the progress.

- Project status: Production-ready
- Current version: 0.9.0 (2025-08-21)
- Changelog: See [CHANGELOG.md](./CHANGELOG.md)

## Features (current)

- **Complete User Authentication System**:
    - User registration with secure password hashing (bcrypt cost factor 12)
    - User login with credential authentication and session management
    - Password validation with 8-character minimum requirement
    - Email format validation with comprehensive regex patterns
    - Duplicate email detection with user-friendly error messages
    - Session-based authentication state tracking
- **CSRF Protection**:
    - Complete protection against Cross-Site Request Forgery attacks
    - CSRF tokens automatically included in all forms
    - Secure CSRF cookie configuration with HttpOnly and Secure flags
- **Route Protection and Access Control**:
    - Authentication middleware for protecting sensitive routes
    - Automatic redirection to login page for unauthenticated users
    - Conditional navigation based on authentication state
    - Cache-Control headers for protected content
- **HTTPS/TLS Server** with production-ready security:
    - Complete TLS implementation with certificate-based encryption
    - Self-signed certificates for development environment
    - Modern cryptographic standards (X25519, P256 curve preferences)
    - Secure session cookies with HTTPS-only flag
    - Connection timeout configurations for production reliability
- **Professional Middleware System**:
    - Request logging with IP tracking and structured logging
    - Comprehensive security headers (CSP, XSS protection, frame options)
    - Panic recovery with graceful error handling
    - Alice middleware chaining for clean composition
    - Multi-layered middleware architecture (public, dynamic, protected routes)
- **Session Management System**:
    - Professional session handling with database storage
    - MySQL-based session store with automatic expiration (12-hour lifetime)
    - Secure session cookie handling with database backing
    - Session token renewal on authentication state changes
    - Protection against session fixation attacks
- **Flash Messaging System**:
    - User feedback with temporary session-based messages
    - Automatic flash message display in templates
    - Success notifications after form submissions
    - Professional styling for user notifications
- **Enhanced Form Handling System**:
    - Professional form processing with validation
    - Custom validation framework with reusable functions
    - Form data preservation on validation errors (sticky forms)
    - Real-time validation error display with field-specific messages
    - Both field and non-field error support
    - Automatic form-to-struct mapping with struct tags
- Structured logging via `log/slog` (startup and error logs)
- Dynamic HTML templates with proper context handling:
    - Base layout template with content blocks (`ui/html/base.tmpl`)
    - Page-specific templates (`ui/html/pages/`)
    - Reusable partial templates (`ui/html/partials/`)
    - Proper context handling in template blocks (using Go's template conventions)
    - Template caching for improved performance
    - Custom template functions (humanDate formatting)
    - Flash message and authentication state integration
- Static assets served from `/static` (`ui/static` for CSS/JS/images)
- Routing with Go 1.22+ pattern-based `ServeMux` (path variables like `{id}`)
- MySQL database integration for persistent snippet and user storage
- Data models with CRUD operations for snippets and users
- **Security Features**:
    - End-to-end TLS encryption for all HTTP traffic
    - Industry-standard password hashing with bcrypt
    - CSRF protection for all state-changing operations
    - Content Security Policy protection
    - Anti-clickjacking headers
    - Content type sniffing protection
    - Secure referrer policy
    - Server-side input validation with comprehensive patterns
    - Length limits and controlled value validation
    - Secure session management with database storage
    - Modern TLS configuration with enhanced cryptographic standards
- Routes (all served over HTTPS with authentication where needed):
    - `/` — home page with latest snippets (public)
    - `/snippet/view/{id}` — view a snippet by numeric ID (public)
    - `/user/signup` — user registration form and processing (public)
    - `/user/login` — user login form and processing (public)
    - `/snippet/create` — create a new snippet (requires authentication)
    - `/user/logout` — user logout (requires authentication)

## Getting started

### Prerequisites

- Go 1.25+ (or compatible)
- MySQL server
- Environment variable `DB_PASSWORD` set with your database password

### Dependencies

The project uses these external libraries:

- `github.com/go-sql-driver/mysql` - MySQL driver for database connectivity
- `github.com/justinas/alice` - HTTP middleware chaining
- `github.com/go-playground/form/v4` - Professional form processing and validation
- `github.com/alexedwards/scs/v2` - Session management framework
- `github.com/alexedwards/scs/mysqlstore` - MySQL-backed session storage
- `github.com/justinas/nosurf` - CSRF protection middleware
- `golang.org/x/crypto` - Cryptographic functions including bcrypt password hashing

### Database Setup

1. Create a MySQL database called `snippetbox`
2. Create the required tables for snippets and users
3. Ensure your MySQL user has appropriate permissions
4. The application uses the DSN format: `web:%s@/snippetbox?parseTime=true` where `%s` is replaced with your password
5. Session data will be automatically stored in the database

### TLS/HTTPS Setup

The application runs exclusively over HTTPS with TLS encryption:

1. **Development Environment**: Self-signed certificates are included in the `tls/` directory
2. **Certificate Files**:
    - `tls/cert.pem` - TLS certificate for localhost
    - `tls/key.pem` - Private key for the certificate
3. **Production**: Replace the development certificates with proper CA-signed certificates

### Run locally

```bash
# From the project root
go run ./...
# Server will start on https://localhost:8080 (note HTTPS)
```

#### Custom address/port

You can change the listen address using the `-addr` flag (defaults to `:8080`):

```bash
go run ./cmd/web -addr=:4000
# Server will start on https://localhost:4000
```

**Important**: The application runs exclusively over HTTPS. Open https://localhost:8080 (or your chosen port) in your
browser. You may need to accept the self-signed certificate warning in your browser for development.

### User Authentication

The application now includes complete user authentication:

1. **Sign Up**: Create a new account at `/user/signup`
    - Requires name, email, and password (minimum 8 characters)
    - Email addresses must be unique
    - Passwords are securely hashed using bcrypt
2. **Log In**: Access your account at `/user/login`
    - Uses email and password authentication
    - Creates secure session for authenticated access
3. **Protected Features**:
    - Creating snippets requires authentication
    - Authenticated users see different navigation options
    - Logout functionality available when logged in

Static files are available under `/static`, for example:

- https://localhost:8080/static/css/main.css
- https://localhost:8080/static/img/logo.png

### Example requests

- Home: `curl -k https://localhost:8080/`
- View snippet: `curl -k https://localhost:8080/snippet/view/123`
- User signup: `curl -k https://localhost:8080/user/signup`
- User login: `curl -k https://localhost:8080/user/login`

Note: Use `-k` flag with curl to accept self-signed certificates in development.

## Development workflow

We maintain a documented history of changes after each section of work.

1. Make changes for the section you are following.
2. Update the `Unreleased` section in [CHANGELOG.md](./CHANGELOG.md) using these categories where applicable:
    - Added, Changed, Deprecated, Removed, Fixed, Security
3. If the section represents a cohesive update, bump the version:
    - Choose the next semantic version (e.g., `0.2.0` for new features).
    - Add the date in `YYYY-MM-DD` format.
4. Commit with a message like:
    - `docs: update changelog for 0.2.0 (2025-08-25)`
5. Push to GitHub so viewers can see progress.

We follow [Semantic Versioning](https://semver.org/) and the [Keep a Changelog](https://keepachangelog.com/) format.

## Versioning policy (summary)

- MAJOR: breaking changes (routes, APIs)
- MINOR: backward-compatible features and improvements
- PATCH: bug fixes and small internal changes

## Project structure (excerpt)

```
cmd/web           # Go entry point and HTTP handlers
  ├─ main.go      # App bootstrap and logging setup (slog)
  ├─ routes.go    # HTTP routes using pattern-based ServeMux
  ├─ handlers.go  # Request handlers
  ├─ templates.go # Template functions and cache
  └─ helpers.go   # Shared helpers (errors, etc.)
internal/models   # Data models and database operations
  ├─ snippets.go  # Snippet model with CRUD operations
  └─ errors.go    # Custom error definitions
ui/html           # Base layout, pages, and partial templates
  ├─ base.tmpl    # Main layout template
  ├─ pages/       # Page-specific templates
  └─ partials/    # Reusable template components
ui/static/css     # Stylesheets
ui/static/js      # JavaScript
ui/static/img     # Images
```

## Template System

Our application uses Go's built-in template package with a structured approach:

### Template Organization

- **Base Layout** (`base.tmpl`): Contains the HTML shell with placeholders for content
- **Page Templates** (`pages/*.tmpl`): Specific content for each route/view
- **Partial Templates** (`partials/*.tmpl`): Reusable components like navigation

### Context Handling

When working with templates, proper context handling is crucial:

- Outside `{{with}}` blocks, use full paths like `.Snippet.Title`
- Inside `{{with .Snippet}}` blocks, the context changes to the Snippet object, so use direct field references like
  `.Title`

### Template Caching

For performance reasons, templates are parsed once at startup and stored in a template cache.

## Roadmap (high level)

- Enhanced HTML templates for server-rendered pages
- Snippet creation form with validation and POST handling
- Basic tests for handlers and routing
- User authentication system
- HTTPS support with automatic certificate management

## Contributing

If you are following along and want to contribute:

- Use conventional commit messages where possible (e.g., `feat:`, `fix:`, `docs:`)
- Update the changelog alongside your changes
- Open a PR describing the section or feature you completed

## License

Add your chosen license here (e.g., MIT). If you include a LICENSE file, link to it from this section.

# Changelog

All notable changes to this project will be documented in this file.

The format is based on "Keep a Changelog" and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned

- Basic tests for handlers and routing

## [0.9.0] - 2025-08-21

### Added

- **Complete User Authentication System** - Full user registration and login functionality
    - `internal/models/users.go` with User model and database operations
    - User registration with secure password hashing using bcrypt (cost factor 12)
    - User login with credential authentication and session management
    - Password validation with minimum 8-character requirement
    - Email format validation with regex pattern matching
    - Duplicate email detection and user-friendly error handling
- **CSRF Protection** - Cross-Site Request Forgery defense system
    - Integration of `github.com/justinas/nosurf` for CSRF token generation
    - CSRF tokens automatically included in all forms via template data
    - Secure CSRF cookie configuration with HttpOnly and Secure flags
    - CSRF middleware applied to all dynamic routes
- **Authentication Middleware** - Route protection and access control
    - `requireAuthentication` middleware for protecting routes requiring login
    - Session-based authentication state checking with `isAuthenticated` helper
    - Automatic redirection to login page for unauthenticated users
    - Cache-Control headers set to prevent caching of protected content
- **Enhanced Validation Framework** - Extended input validation capabilities
    - Email regex validation with comprehensive pattern matching
    - Minimum character length validation (`MinChars`)
    - Pattern matching validation (`Matches`) with regex support
    - Non-field error support for general form validation messages
    - Enhanced Validator struct with both field and non-field error handling
- **User Authentication Templates** - Complete UI for user management
    - `signup.tmpl` template with name, email, and password fields
    - `login.tmpl` template with email/password authentication
    - Form validation error display with field-specific messaging
    - Non-field error display for authentication failures
    - CSRF token integration in all authentication forms
- **Session Security Enhancements** - Advanced session management features
    - Session token renewal on login/logout for security
    - Authenticated user ID storage in session data
    - Session-based authentication state tracking
    - Secure logout with session data cleanup and token renewal

### Changed

- **Application Architecture** - Enhanced with user management capabilities
    - Added `users *models.UserModel` to application struct
    - Updated template data structure with `IsAuthenticated` and `CSRFToken` fields
    - Enhanced `newTemplateData` helper to automatically populate auth status and CSRF tokens
- **Route Protection** - Access control implementation
    - Snippet creation routes now require authentication
    - User logout route protected and accessible only when authenticated
    - Clear separation between public and protected route groups
- **Navigation System** - Dynamic UI based on authentication state
    - Conditional navigation menu showing different options for authenticated/unauthenticated users
    - "Create Snippet" link only visible to authenticated users
    - Login/Signup links for unauthenticated users
    - Logout form with CSRF protection for authenticated users
- **Error Handling** - Enhanced validation and authentication error management
    - New authentication-specific error types: `ErrInvalidCredentials`, `ErrDuplicateEmail`
    - Comprehensive form validation with both field and non-field errors
    - User-friendly error messages for authentication failures
- **Middleware Architecture** - Multi-layered request processing
    - Dynamic routes now use both session and CSRF protection middleware
    - Protected routes use additional authentication middleware layer
    - Enhanced middleware composition with Alice chaining

### Security

- **Password Security** - Industry-standard password protection
    - Bcrypt password hashing with high cost factor (12)
    - Secure password storage in database with hashed_password field
    - Password comparison using bcrypt's constant-time comparison
- **Session Security** - Enhanced session management
    - Session token renewal on authentication state changes
    - Secure session data handling for user identification
    - Protection against session fixation attacks
- **CSRF Protection** - Complete protection against cross-site attacks
    - CSRF tokens required for all state-changing operations
    - Secure CSRF cookie configuration
    - Integration with all forms requiring protection

### Dependencies

- **New Security Libraries** - Professional authentication and security tools
    - `github.com/justinas/nosurf v1.2.0` - CSRF protection middleware
    - `golang.org/x/crypto v0.41.0` - Cryptographic functions including bcrypt

### Infrastructure

- **Database Schema Extensions** - User management tables
    - Users table with id, name, email, hashed_password, and created fields
    - Email uniqueness constraint for preventing duplicate accounts
    - MySQL integration with existing database structure

## [0.8.0] - 2025-08-21

### Added

- **HTTPS/TLS Support** - Complete secure server implementation
    - TLS server configuration with `ListenAndServeTLS()`
    - Self-signed certificate generation for development (`./tls/cert.pem`, `./tls/key.pem`)
    - Custom TLS configuration with modern curve preferences (X25519, P256)
    - Secure session cookies with `Secure: true` flag
- **Enhanced Server Configuration** - Professional server setup
    - Structured `http.Server` configuration with custom settings
    - Connection timeout configurations:
        - `IdleTimeout: time.Minute` - Connection idle timeout
        - `ReadTimeout: 5 * time.Second` - Request read timeout
        - `WriteTimeout: 10 * time.Second` - Response write timeout
    - Custom error logging integration with slog
    - TLS certificate file management and organization
- **Security Enhancements** - Production-ready security features
    - HTTPS-only session cookies for secure session management
    - Modern TLS curve preferences for enhanced encryption
    - Certificate-based encryption for all HTTP traffic
    - Secure localhost development environment setup

### Changed

- **Server Architecture** - Enhanced from basic HTTP to secure HTTPS
    - Migrated from `http.ListenAndServe()` to `srv.ListenAndServeTLS()`
    - Added comprehensive server configuration structure
    - Enhanced session security with HTTPS-only cookies
    - Updated logging to use structured error logging with slog integration
- **Development Environment** - HTTPS-first development setup
    - Local development now uses HTTPS with self-signed certificates
    - Enhanced security posture for development and testing
    - Certificate management for development workflows

### Security

- **Transport Layer Security** - End-to-end encryption
    - All HTTP traffic now encrypted with TLS
    - Modern cryptographic standards with elliptic curve preferences
    - Secure session cookie handling prevents session hijacking
    - Certificate-based authentication and encryption

### Infrastructure

- **TLS Certificate Management** - Organized certificate structure
    - `tls/` directory for certificate and key storage
    - Self-signed certificates for development environment
    - Proper file organization for production certificate deployment

## [0.7.0] - 2025-08-21

### Added

- **Session Management System** - Complete session handling with database storage
    - Integration of `github.com/alexedwards/scs/v2` for professional session management
    - MySQL-based session storage using `github.com/alexedwards/scs/mysqlstore`
    - 12-hour session lifetime with automatic expiration
    - Session middleware integration with Alice middleware chains
- **Flash Messaging System** - User feedback with temporary messages
    - Flash message support for user notifications and feedback
    - Session-based flash message storage with automatic cleanup
    - Template integration for flash message display
    - Success message display after snippet creation
- **Enhanced User Experience** - Improved feedback and interaction
    - Flash message styling with professional appearance
    - Automatic flash message display in base template
    - Context-aware message handling with session integration
    - User feedback after form submissions
- **Advanced Middleware Architecture** - Sophisticated request processing
    - Dynamic middleware chain for session-enabled routes
    - Separation of static and dynamic route handling
    - Session middleware (`LoadAndSave`) integration with existing middleware
    - Clean middleware composition with Alice chaining

### Changed

- **Application Structure** - Enhanced with session capabilities
    - Added `sessionManager *scs.SessionManager` to application struct
    - Session manager initialization in main.go bootstrap
    - Updated imports to include session management libraries
- **Template System** - Flash message integration
    - Added `Flash string` field to `templateData` struct
    - Enhanced `newTemplateData()` helper to auto-populate flash messages
    - Base template updated with flash message display block
- **Route Architecture** - Session-aware routing
    - Implemented dynamic middleware chain for session-enabled routes
    - All dynamic routes now use session middleware
    - Maintained static file serving without session overhead
- **Handler Enhancement** - User feedback integration
    - Updated `snippetCreatePost` handler to set success flash messages
    - Enhanced `snippetView` handler with flash message support (prepared but not active)
    - Improved user feedback workflow after form submissions

### Security

- **Session Security** - Secure session management
    - Database-backed session storage for security and scalability
    - Automatic session expiration (12-hour lifetime)
    - Secure session cookie handling
    - Session data isolated from client-side storage

### Dependencies

- **New Libraries** - Professional session management
    - `github.com/alexedwards/scs/v2 v2.9.0` - Core session management
    - `github.com/alexedwards/scs/mysqlstore` - MySQL session store

## [0.6.0] - 2025-08-21

### Added

- **Complete Form Handling System** - Professional form processing architecture
    - `snippetCreateForm` struct with embedded validator for form data handling
    - Form field validation with custom error messages
    - Form data preservation on validation errors (sticky forms)
    - Proper form encoding/decoding with struct tags
- **Validation Framework** - Comprehensive input validation system
    - New `internal/validator` package with reusable validation functions
    - `Validator` struct with field error mapping and validation state tracking
    - Validation helper functions:
        - `NotBlank()` - ensures fields are not empty
        - `MaxChars()` - enforces character limits with UTF-8 support
        - `PermittedValues()` - validates against allowed values using generics
    - Embedded validator pattern for clean form struct integration
- **Form Processing Library Integration** - Professional form handling
    - Added `github.com/go-playground/form/v4` dependency for form decoding
    - `decodePostForm()` helper method for automatic form-to-struct mapping
    - Proper error handling for form decoding with panic recovery
    - Form decoder initialization in application bootstrap
- **Create Snippet Form** - Complete user input interface
    - New `create.tmpl` template with full form implementation
    - Form fields: title (text), content (textarea), expires (radio buttons)
    - Real-time validation error display with field-specific messages
    - Form value preservation on validation errors
    - Proper form submission handling with POST method
- **Enhanced UI/UX** - Professional form styling and navigation
    - Comprehensive form CSS styling with error state handling
    - Error styling with red borders and bold error messages
    - Navigation integration with "Create Snippet" link
    - Responsive form layout with consistent spacing
    - Radio button styling for expiration options (1 day, 1 week, 1 year)

### Changed

- **Handler Architecture** - Enhanced request processing
    - `snippetCreate` GET handler now renders proper form template with defaults
    - `snippetCreatePost` POST handler implements full validation workflow
    - Template data structure updated with generic `Form any` field
    - Integration of validation workflow with template rendering
- **Application Structure** - Form processing capabilities
    - Added `formDecoder *form.Decoder` to application struct
    - Form decoder initialization in main.go bootstrap
    - Updated imports to include form processing and validation packages
- **Template System** - Form-aware template rendering
    - Enhanced `templateData` struct to support any form type
    - Template integration with validation error display
    - Conditional rendering based on validation state
- **Error Handling** - Improved form error processing
    - HTTP 422 Unprocessable Entity status for validation errors
    - Graceful form re-rendering on validation failures
    - Structured error display in templates

### Security

- **Input Validation** - Defense against malicious input
    - Server-side validation for all form fields
    - Length limits on text inputs to prevent buffer attacks
    - Controlled value validation for restricted fields
    - Proper form parsing with error handling

## [0.5.0] - 2025-08-20

### Added

- **HTTP Middleware System** - Complete middleware architecture for request processing
    - `middleware.go` with three core middleware functions:
        - `commonHeaders()` - Security headers and server identification
        - `logRequest()` - Structured request logging with IP, method, URI, and protocol
        - `recoverPanic()` - Panic recovery with graceful error handling
- **Security Headers Implementation** - Comprehensive security header configuration:
    - Content Security Policy (CSP) with font and style source restrictions
    - Referrer Policy set to "origin-when-cross-origin"
    - X-Content-Type-Options: "nosniff"
    - X-Frame-Options: "deny"
    - X-XSS-Protection: "0" (modern approach)
    - Custom Server header set to "Go"
- **Alice Middleware Library Integration** - Professional middleware chaining
    - Added `github.com/justinas/alice v1.2.0` dependency
    - Implemented middleware chain pattern in routes for clean composition
    - Standard middleware chain: `recoverPanic` → `logRequest` → `commonHeaders`
- **Enhanced Request Logging** - Detailed request tracking
    - IP address logging for security and analytics
    - HTTP protocol version tracking
    - Method and URI logging for debugging
    - Integration with existing slog structured logging

### Changed

- **Routes Architecture** - Updated routing system to use middleware chains
    - Refactored `routes.go` to implement Alice middleware chaining
    - All routes now pass through the standard middleware chain
    - Improved separation of concerns between routing and middleware
- **Error Handling** - Enhanced panic recovery and error reporting
    - Connection close header set on panic recovery
    - Graceful degradation on server errors
    - Consistent error logging through middleware chain

### Security

- **Multiple Security Headers** - Defense against common web vulnerabilities
    - CSP protection against XSS and injection attacks
    - Frame options to prevent clickjacking
    - Content type sniffing protection
    - Referrer policy for privacy protection

## [0.4.1] - 2025-08-20

### Added

- Custom template function `humanDate` for formatting time values in a user-friendly format
- Buffer-based template rendering to improve error handling and performance
- Template data helper function `newTemplateData` that automatically includes the current year
- Comprehensive documentation for dynamic HTML templates system explaining the structure and context handling

### Fixed

- Template error in view.tmpl when accessing individual snippets - corrected context handling within the {{with
  .Snippet}} block by using direct field references (.Title, .ID, etc.) instead of redundant path notation (
  .Snippet.Title)
- Improved template context handling to follow Go's standard template conventions
- Enhanced error handling in template rendering to provide clearer error messages

### Changed

- Optimized template execution with a buffered approach to catch errors before writing to the response
- Enhanced README documentation with detailed template system architecture
- Updated project structure documentation to highlight the template organization
- Added detailed explanations of Go template context handling in documentation

## [0.4.0] - 2025-08-19

### Added

- MySQL database integration for persistent snippet storage
- Database connection setup with environment-based password configuration
- `internal/models` package with data models and database operations:
    - `Snippet` struct representing the data model
    - `SnippetModel` for database operations (Insert, Get, Latest)
    - Custom error handling with `ErrNoRecord`
- Command-line flag for database connection string (`-dsn`)
- Database connection pooling and proper resource cleanup
- Database-powered snippet routes:
    - Home page now displays latest snippets from database
    - View snippet fetches data from database by ID
    - Create snippet endpoint stores data in database

### Changed

- Updated application structure to support dependency injection of database
- Handlers now use the model layer to access data instead of hardcoded responses
- Added database connection details to documentation

## [0.3.0] - 2025-08-19

### Added

- Structured application logging using `log/slog` (startup and error logs)
- Dedicated error handling helpers in `helpers.go`:
    - `serverError` for internal 500 errors with detailed logging
    - `clientError` for general HTTP error responses

### Changed

- Refactored project file structure in `cmd/web` to separate concerns:
    - Introduced `routes.go` for HTTP route registrations
    - Introduced `helpers.go` for shared error/helper functions
- Upgraded Go version requirement to 1.25 (updated `go.mod` and README prerequisites)
- Documentation: Updated README to include structured logging, refined project structure, and Go 1.25 prerequisite

## [0.2.0] - 2025-08-19

### Added

- Server-side HTML template rendering for the home page (base layout, nav partial, home page)
- Static file serving from `/static` (CSS, JS, images); added favicon and logo assets
- GET and POST handlers for `/snippet/create` with basic responses
- Basic UI scaffolding: `ui/static/css/main.css` and `ui/static/js/main.js`
- Route for viewing specific snippets with ID parameter (`/snippet/view/{id}`)

### Changed

- Home route now renders templates instead of plain text
- Routing now uses Go 1.22 pattern-based `ServeMux` with path parameters (e.g., `{id}`)
- Documentation: Expanded README with details on templates and static assets, browser usage, and project structure (
  2025-08-18)

## [0.1.0] - 2025-08-18

### Added

- Initial project structure and Go module setup
- Basic HTTP server with `net/http`
- Command-line flag for custom address/port configuration
- Simple handler functions for home, snippet view, and snippet creation
- Project documentation in README.md with setup and usage instructions

---

### How we version

- Patch (x.y.Z): Bug fixes and small internal changes that do not add features
- Minor (x.Y.z): Backwards-compatible feature additions and improvements
- Major (X.y.z): Breaking changes in API, routes, or behavior

### How to update this changelog after each section

1. Add your changes under the `Unreleased` section using the categories: `Added`, `Changed`, `Deprecated`, `Removed`,
   `Fixed`, `Security`.
2. When you are ready to tag a version:
    - Decide the next version number (e.g., 0.2.0 for a new feature set).
    - Replace `Unreleased` with a new version heading including the date, and create a fresh empty `Unreleased` section
      above it.
3. Commit with a message like: `docs: update changelog for 0.2.0 (2025-08-25)`.

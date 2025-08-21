# Changelog

All notable changes to this project will be documented in this file.

The format is based on "Keep a Changelog" and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned

- Basic tests for handlers and routing
- User authentication system
- HTTPS support with automatic certificate management

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

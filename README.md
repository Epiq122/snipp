# Snippet

A simple Go web application for creating and viewing text snippets. This repository is being developed incrementally,
with each section documented and versioned so that viewers can follow the progress.

- Project status: Active development
- Current version: 0.7.0 (2025-08-21)
- Changelog: See [CHANGELOG.md](./CHANGELOG.md)

## Features (current)

- HTTP server using `net/http`
- **Professional Middleware System**:
    - Request logging with IP tracking and structured logging
    - Comprehensive security headers (CSP, XSS protection, frame options)
    - Panic recovery with graceful error handling
    - Alice middleware chaining for clean composition
    - Dynamic middleware chains for session-enabled routes
- **Session Management System**:
    - Professional session handling with database storage
    - MySQL-based session store with automatic expiration (12-hour lifetime)
    - Secure session cookie handling with database backing
    - Session middleware integration with existing middleware chains
- **Flash Messaging System**:
    - User feedback with temporary session-based messages
    - Automatic flash message display in templates
    - Success notifications after form submissions
    - Professional styling for user notifications
- **Complete Form Handling System**:
    - Professional form processing with validation
    - Custom validation framework with reusable functions
    - Form data preservation on validation errors (sticky forms)
    - Real-time validation error display with field-specific messages
    - Automatic form-to-struct mapping with struct tags
- Structured logging via `log/slog` (startup and error logs)
- Dynamic HTML templates with proper context handling:
    - Base layout template with content blocks (`ui/html/base.tmpl`)
    - Page-specific templates (`ui/html/pages/`)
    - Reusable partial templates (`ui/html/partials/`)
    - Proper context handling in template blocks (using Go's template conventions)
    - Template caching for improved performance
    - Custom template functions (humanDate formatting)
    - Flash message integration in base template
- Static assets served from `/static` (`ui/static` for CSS/JS/images)
- Routing with Go 1.22+ pattern-based `ServeMux` (path variables like `{id}`)
- MySQL database integration for persistent snippet storage
- Data models with CRUD operations for snippets
- **Security Features**:
    - Content Security Policy protection
    - Anti-clickjacking headers
    - Content type sniffing protection
    - Secure referrer policy
    - Server-side input validation
    - Length limits and controlled value validation
    - Secure session management with database storage
- Routes
    - `/` — home page with latest snippets
    - `/snippet/view/{id}` — view a snippet by numeric ID (with flash message support)
    - `/snippet/create` — create a new snippet (GET: form, POST: processing with validation and success feedback)

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

### Database Setup

1. Create a MySQL database called `snippetbox`
2. Ensure your MySQL user has appropriate permissions
3. The application uses the DSN format: `web:%s@/snippetbox?parseTime=true` where `%s` is replaced with your password
4. Session data will be automatically stored in the database

### Run locally

```bash
# From the project root
go run ./...
# Server will start on http://localhost:8080
```

#### Custom address/port

You can change the listen address using the `-addr` flag (defaults to `:8080`):

```bash
go run ./cmd/web -addr=:4000
# Server will start on http://localhost:4000
```

Then open http://localhost:8080 (or your chosen port) in your browser to view the templated home page.

Static files are available under `/static`, for example:

- http://localhost:8080/static/css/main.css
- http://localhost:8080/static/img/logo.png

### Example requests

- Home: `curl http://localhost:8080/`
- View snippet: `curl http://localhost:8080/snippet/view/123`
- Create (GET placeholder): `curl http://localhost:8080/snippet/create`
- Create (POST placeholder): `curl -i -X POST http://localhost:8080/snippet/create`

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

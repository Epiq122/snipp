# Changelog

All notable changes to this project will be documented in this file.

The format is based on "Keep a Changelog" and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned

- Add persistent storage for snippets (database layer)
- Add create form with validation and POST handling
- Basic tests for handlers and routing
- User authentication system
- HTTPS support with automatic certificate management

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

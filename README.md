# Snippet

A simple Go web application for creating and viewing text snippets. This repository is being developed incrementally, with each section documented and versioned so that viewers can follow the progress.

- Project status: Early development
- Current version: 0.2.0 (2025-08-19)
- Changelog: See [CHANGELOG.md](./CHANGELOG.md)

## Features (current)
- HTTP server using `net/http`
- Server-side HTML templates for the home page (`ui/html`)
- Static assets served from `/static` (`ui/static` for CSS/JS/images)
- Routes
  - `/` — home page rendered via templates
  - `/snippet/view/{id}` — view a snippet by numeric ID
  - `/snippet/create` — placeholder for creating a snippet (GET and POST)

## Getting started

### Prerequisites
- Go 1.22+ (or compatible)

### Run locally
```bash
# From the project root
go run ./...
# Server will start on http://localhost:8080
```

Then open http://localhost:8080 in your browser to view the templated home page.

Static files are available under `/static`, for example:
- http://localhost:8080/static/css/main.css
- http://localhost:8080/static/img/logo.png

### Example requests
- Home: `curl http://localhost:8080/`
- View snippet: `curl http://localhost:8080/snippet/view/123`
- Create (placeholder): `curl http://localhost:8080/snippet/create`

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
ui/html           # Base layout, pages, and partial templates
ui/static/css     # Stylesheets
ui/static/js      # JavaScript
ui/static/img     # Images
```

## Roadmap (high level)
- Persistent storage for snippets (database)
- Enhanced HTML templates for server-rendered pages
- Snippet creation form with validation and POST handling
- Basic tests for handlers and routing

## Contributing
If you are following along and want to contribute:
- Use conventional commit messages where possible (e.g., `feat:`, `fix:`, `docs:`)
- Update the changelog alongside your changes
- Open a PR describing the section or feature you completed

## License
Add your chosen license here (e.g., MIT). If you include a LICENSE file, link to it from this section.

# Changelog

All notable changes to this project will be documented in this file.

The format is based on "Keep a Changelog" and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
- Planned: Add persistent storage for snippets (database layer)
- Planned: Implement HTML templates for views
- Planned: Add create form with validation and POST handling
- Planned: Basic tests for handlers and routing

## [0.1.0] - 2025-08-18
### Added
- Initial Go module and project scaffold
- HTTP server with `net/http` listening on `:8080`
- Routes:
  - `/` (home): returns a greeting
  - `/snippet/view/{id}`: displays a snippet ID parsed from the URL
  - `/snippet/create`: placeholder endpoint indicating snippet creation
- Basic server startup logging

### Notes
- This is an early development snapshot suitable for following along with incremental sections. No persistence or templates yet.

---

### How we version
- Patch (x.y.Z): Bug fixes and small internal changes that do not add features
- Minor (x.Y.z): Backwards-compatible feature additions and improvements
- Major (X.y.z): Breaking changes in API, routes, or behavior

### How to update this changelog after each section
1. Add your changes under the `Unreleased` section using the categories: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, `Security`.
2. When you are ready to tag a version:
   - Decide the next version number (e.g., 0.2.0 for a new feature set).
   - Replace `Unreleased` with a new version heading including the date, and create a fresh empty `Unreleased` section above it.
3. Commit with a message like: `docs: update changelog for 0.2.0 (2025-08-25)`.

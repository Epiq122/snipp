# Go Web Development Project Guide

This document serves as a comprehensive guide for Go web application development, based on the Snippet project. It
outlines key components, best practices, and implementation patterns that can be reused in future projects.

## Project Structure

A well-organized project structure helps maintain code quality and promotes separation of concerns:

```
cmd/web/               # Application entry point and web components
  ├─ main.go           # App bootstrap, configuration, and dependency injection
  ├─ routes.go         # HTTP route definitions with middleware chains
  ├─ handlers.go       # HTTP request handlers
  ├─ middleware.go     # HTTP middleware functions
  ├─ templates.go      # Template functions and cache
  └─ helpers.go        # Shared helper functions
internal/models/       # Data models and database operations
  ├─ snippets.go       # Core model definitions with CRUD operations
  └─ errors.go         # Custom error types
ui/                    # User interface components
  ├─ html/             # HTML templates
  │   ├─ base.tmpl     # Base layout template
  │   ├─ pages/        # Page-specific templates
  │   └─ partials/     # Reusable template components
  └─ static/           # Static assets
      ├─ css/          # Stylesheets
      ├─ js/           # JavaScript files
      └─ img/          # Images and icons
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
logger        *slog.Logger
snippets      *models.SnippetModel
templateCache map[string]*template.Template
}

// Environment variable handling for sensitive data
password := os.Getenv("DB_PASSWORD")
if password == "" {
logger.Error("DB_PASSWORD environment variable not set")
os.Exit(1)
}
```

### 2. Database Integration Pattern

```go
// Database connection with proper error handling
func openDB(dsn string) (*sql.DB, error) {
db, err := sql.Open("mysql", dsn)
if err != nil {
return nil, err
}

err = db.Ping()
if err != nil {
db.Close()
return nil, err
}
return db, nil
}

// Model structure for database operations
type SnippetModel struct {
DB *sql.DB
}

// CRUD operations with proper error handling
func (m *SnippetModel) Get(id int) (Snippet, error) {
stmt := `SELECT id, title, content, created, expires FROM snippets
            WHERE expires > UTC_TIMESTAMP() AND id = ?`

row := m.DB.QueryRow(stmt, id)

var s Snippet
err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
if err != nil {
if errors.Is(err, sql.ErrNoRows) {
return Snippet{}, ErrNoRecord
}
return Snippet{}, err
}
return s, nil
}
```

### 3. HTTP Middleware System

Professional middleware architecture using Alice for chaining:

```go
// Security headers middleware
func commonHeaders(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Security-Policy",
"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "deny")
w.Header().Set("X-XSS-Protection", "0")
w.Header().Set("Server", "Go")
next.ServeHTTP(w, r)
})
}

// Request logging middleware
func (app *application) logRequest(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
var (
ip = r.RemoteAddr
proto = r.Proto
method = r.Method
uri = r.URL.RequestURI()
)
app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
next.ServeHTTP(w, r)
})
}

// Panic recovery middleware
func (app *application) recoverPanic(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
defer func () {
if err := recover(); err != nil {
w.Header().Set("Connection", "close")
app.serverError(w, r, fmt.Errorf("%s", err))
}
}()
next.ServeHTTP(w, r)
})
}

// Middleware chaining with Alice
func (app *application) routes() http.Handler {
mux := http.NewServeMux()
// ... route definitions ...

standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
return standard.Then(mux)
}
```

### 4. Template System Architecture

```go
// Template data structure
type templateData struct {
CurrentYear int
Snippet     models.Snippet
Snippets    []models.Snippet
}

// Template cache for performance
func newTemplateCache() (map[string]*template.Template, error) {
cache := map[string]*template.Template{}
pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
if err != nil {
return nil, err
}

for _, page := range pages {
name := filepath.Base(page)

// Register custom functions before parsing
ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
if err != nil {
return nil, err
}

// Add partials and page templates
ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
if err != nil {
return nil, err
}

ts, err = ts.ParseFiles(page)
if err != nil {
return nil, err
}

cache[name] = ts
}
return cache, nil
}

// Custom template functions
var functions = template.FuncMap{
"humanDate": humanDate,
}

func humanDate(t time.Time) string {
return t.Format("02 Jan 2006 at 15:04")
}
```

### 5. Error Handling Patterns

```go
// Structured error handling with logging
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
var (
method = r.Method
uri = r.URL.RequestURI()
)
app.logger.Error(err.Error(), "method", method, "url", uri)
http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Buffer-based template rendering for error catching
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
ts, ok := app.templateCache[page]
if !ok {
err := fmt.Errorf("template %s not found", page)
app.serverError(w, r, err)
return
}

buf := new(bytes.Buffer)
err := ts.ExecuteTemplate(buf, "base", data)
if err != nil {
app.serverError(w, r, err)
return
}

w.WriteHeader(status)
buf.WriteTo(w)
}
```

### 6. Handler Patterns

```go
// Standard handler pattern with error handling
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
id, err := strconv.Atoi(r.PathValue("id"))
if err != nil {
http.NotFound(w, r)
return
}

snippet, err := app.snippets.Get(id)
if err != nil {
if errors.Is(err, models.ErrNoRecord) {
http.NotFound(w, r)
return
} else {
app.serverError(w, r, err)
}
return
}

data := app.newTemplateData(r)
data.Snippet = snippet
app.render(w, r, http.StatusOK, "view.tmpl", data)
}
```

## Security Best Practices Implemented

### 1. Security Headers

- **Content Security Policy (CSP)**: Prevents XSS attacks by controlling resource loading
- **X-Frame-Options**: Prevents clickjacking attacks
- **X-Content-Type-Options**: Prevents MIME type sniffing
- **Referrer-Policy**: Controls referrer information sent with requests
- **X-XSS-Protection**: Set to 0 following modern security practices

### 2. Database Security

- Environment variable for database passwords
- Prepared statements to prevent SQL injection
- Proper connection string formatting

### 3. Error Handling

- Generic error responses to avoid information disclosure
- Detailed logging for debugging without exposing sensitive data
- Panic recovery to prevent application crashes

## Development Workflow Best Practices

### 1. Project Versioning

- Follow Semantic Versioning (SemVer)
- Maintain detailed changelog with categorized changes
- Document all major features and fixes

### 2. Code Organization

- Clear separation of concerns (handlers, models, templates)
- Dependency injection pattern for testability
- Consistent error handling throughout the application

### 3. Documentation Standards

- Document all public functions and types
- Maintain comprehensive README with setup instructions
- Keep changelog updated with each release
- Include examples and usage patterns

## Dependencies and External Libraries

### Core Dependencies

- `github.com/go-sql-driver/mysql` - MySQL database driver
- `github.com/justinas/alice` - HTTP middleware chaining

### Standard Library Usage

- `net/http` - HTTP server and routing
- `html/template` - Template rendering
- `log/slog` - Structured logging
- `database/sql` - Database interface
- `flag` - Command-line argument parsing

## Testing Considerations

### Areas to Test

- Handler functions with various inputs
- Database model operations
- Template rendering
- Middleware functionality
- Error handling scenarios

### Testing Structure

```
cmd/web/
  ├─ handlers_test.go
  ├─ middleware_test.go
  └─ templates_test.go
internal/models/
  └─ snippets_test.go
```

## Deployment Considerations

### Environment Variables

- `DB_PASSWORD` - Database password
- Consider additional config for production (log levels, timeouts, etc.)

### Production Readiness Checklist

- [ ] HTTPS configuration
- [ ] Database connection pooling tuned for load
- [ ] Logging configured for production environment
- [ ] Security headers properly configured
- [ ] Static asset serving optimized
- [ ] Database migrations handled
- [ ] Health check endpoints
- [ ] Graceful shutdown handling

## Future Enhancements

### Planned Features

- Form handling with validation
- User authentication and sessions
- CSRF protection
- Rate limiting
- API endpoints with JSON responses
- File upload handling
- Email notifications
- Admin interface

### Scalability Considerations

- Database connection pooling optimization
- Caching layer (Redis/Memcached)
- Load balancer compatibility
- Horizontal scaling patterns
- Monitoring and metrics

## Lessons Learned

### What Worked Well

- Alice middleware chaining provides clean, composable architecture
- Template caching significantly improves performance
- Structured logging makes debugging easier
- Buffer-based template rendering catches errors before sending responses
- Dependency injection makes the application testable

### Common Pitfalls to Avoid

- Template context confusion (use correct field references within `{{with}}` blocks)
- Missing error handling in database operations
- Forgetting to set proper security headers
- Not using prepared statements for SQL queries
- Hardcoding configuration values instead of using flags/environment variables

## Conclusion

This guide represents a solid foundation for Go web applications, incorporating:

- Professional middleware architecture
- Secure by default configurations
- Clean code organization
- Comprehensive error handling
- Performance optimizations
- Development best practices

The patterns demonstrated here can be adapted and extended for more complex applications while maintaining code quality
and security standards.

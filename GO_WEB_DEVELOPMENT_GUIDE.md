# Go Web Development Project Guide

This document serves as a comprehensive guide for Go web application development, based on the Snippet project. It outlines key components, best practices, and implementation patterns that can be reused in future projects.

## Project Structure

A well-organized project structure helps maintain code quality and promotes separation of concerns:

```
cmd/web/               # Application entry point and web components
  ├─ main.go           # App bootstrap, configuration, and dependency injection
  ├─ routes.go         # HTTP route definitions
  ├─ handlers.go       # HTTP request handlers
  ├─ templates.go      # Template functions and cache
  └─ helpers.go        # Shared helper functions
internal/models/       # Data models and database operations
  ├─ models.go         # Core model definitions
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

### 1. Application Configuration

```go
// Command-line flag parsing for configuration
flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
flag.StringVar(&cfg.dsn, "dsn", fmt.Sprintf("web:%s@/snippetbox?parseTime=true", os.Getenv("DB_PASSWORD")), "MySQL data source name")
flag.Parse()

// Application struct for dependency injection
type application struct {
    logger    *slog.Logger
    snippets  *models.SnippetModel
    templateCache map[string]*template.Template
}
```

### 2. Database Integration

```go
// Open a connection pool
db, err := openDB(cfg.dsn)
if err != nil {
    logger.Error(err.Error())
    os.Exit(1)
}
defer db.Close()

// Database helper function
func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
```

### 3. Routing with Go 1.22+ Pattern Matching

```go
// Initialize a new ServeMux instance
mux := http.NewServeMux()

// Register routes with pattern matching
mux.HandleFunc("GET /", app.home)
mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
mux.HandleFunc("GET /snippet/create", app.snippetCreate)
mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

// Serve static files
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
```

### 4. Template System

#### Template Setup and Caching

```go
// Create template cache at application startup
templateCache, err := newTemplateCache()
if err != nil {
    logger.Error(err.Error())
    os.Exit(1)
}

func newTemplateCache() (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}
    
    // Find all page templates
    pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
    if err != nil {
        return nil, err
    }
    
    for _, page := range pages {
        name := filepath.Base(page)
        
        // Parse base template with custom functions
        ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
        if err != nil {
            return nil, err
        }
        
        // Add partials
        ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
        if err != nil {
            return nil, err
        }
        
        // Add the page template
        ts, err = ts.ParseFiles(page)
        if err != nil {
            return nil, err
        }
        
        cache[name] = ts
    }
    
    return cache, nil
}
```

#### Template Data Structure

```go
// Shared data structure for templates
type templateData struct {
    CurrentYear int
    Snippet     models.Snippet
    Snippets    []models.Snippet
}

// Helper to create template data with common fields
func (app *application) newTemplateData(r *http.Request) templateData {
    return templateData{
        CurrentYear: time.Now().Year(),
    }
}
```

#### Custom Template Functions

```go
// Format time in a human-readable format
func humanDate(t time.Time) string {
    return t.Format("02 Jan 2006 at 15:04")
}

// Register custom functions
var functions = template.FuncMap{
    "humanDate": humanDate,
}
```

#### Buffered Template Rendering

```go
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
    // Get template from cache
    ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("template %s not found", page)
        app.serverError(w, r, err)
        return
    }
    
    // Use a buffer for template execution to catch errors before writing response
    buf := new(bytes.Buffer)
    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        app.serverError(w, r, err)
        return
    }
    
    // Set status code and write buffered content
    w.WriteHeader(status)
    buf.WriteTo(w)
}
```

### 5. Error Handling

```go
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
    var (
        method = r.Method
        uri    = r.URL.RequestURI()
    )
    
    // Log detailed error information
    app.logger.Error(err.Error(), "method", method, "url", uri)
    
    // Send generic error response to user
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}
```

### 6. Model Operations

```go
// Example model with database operations
type SnippetModel struct {
    DB *sql.DB
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
    stmt := `SELECT id, title, content, created, expires FROM snippets
             WHERE id = ? AND expires > UTC_TIMESTAMP()`
             
    var s Snippet
    err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return Snippet{}, ErrNoRecord
        }
        return Snippet{}, err
    }
    
    return s, nil
}
```

### 7. Handler Implementation Pattern

```go
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    // Parse path parameter
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }
    
    // Get data from model
    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.clientError(w, http.StatusNotFound)
        } else {
            app.serverError(w, r, err)
        }
        return
    }
    
    // Initialize template data
    data := app.newTemplateData(r)
    data.Snippet = snippet
    
    // Render template
    app.render(w, r, http.StatusOK, "view.tmpl", data)
}
```

## Template Context Handling

When working with Go's template system, proper context handling is essential:

1. **Outside `{{with}}` blocks**: Use full dot notation paths to access data
   ```
   {{.Snippet.Title}}
   ```

2. **Inside `{{with .Snippet}}` blocks**: The context (dot) changes to be the Snippet value, so use direct field references
   ```
   {{with .Snippet}}
       <h1>{{.Title}}</h1>   <!-- Not {{.Snippet.Title}} -->
   {{end}}
   ```

3. **Range blocks**: Similarly, the dot context changes inside range loops
   ```
   {{range .Snippets}}
       <div>{{.Title}}</div>   <!-- Each individual snippet -->
   {{end}}
   ```

## Structured Logging

Use `log/slog` for structured, level-based logging:

```go
// Initialize structured logger
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// Log with additional context fields
logger.Info("server starting", "addr", cfg.addr)
logger.Error(err.Error(), "method", method, "url", uri)
```

## Version Management

Follow semantic versioning for releases:

- **Patch (0.0.x)**: Bug fixes and small internal changes
- **Minor (0.x.0)**: New features and non-breaking improvements
- **Major (x.0.0)**: Breaking changes to APIs or behavior

Document all changes in a CHANGELOG.md file using "Keep a Changelog" format.

## Common Gotchas and Solutions

1. **Template Context Issues**: Remember that inside `{{with}}` blocks, the dot context changes.

2. **Database Connections**: Always check connections with `db.Ping()` and implement proper connection pooling.

3. **Error Handling**: Implement consistent error handling patterns and avoid exposing internal errors to users.

4. **Path Parameters**: Use proper type conversion and validation for URL parameters.

5. **Static File Serving**: Use `http.StripPrefix` when serving static files from a subdirectory.

6. **Template Errors**: Use a buffer to execute templates first to catch errors before writing to the response.

7. **HTTP Methods**: Explicitly specify HTTP methods in route handlers to avoid security issues.

---

This guide is based on the Snippet project development as of version 0.4.1 (2025-08-20).

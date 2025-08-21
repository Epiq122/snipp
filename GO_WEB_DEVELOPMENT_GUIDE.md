# Go Web Development Project Guide

This document serves as a comprehensive guide for Go web application development, based on the Snippet project. It
outlines key components, best practices, and implementation patterns that can be reused in future projects.

## Project Structure

A well-organized project structure helps maintain code quality and promotes separation of concerns:

```
cmd/web/               # Application entry point and web components
  ├─ main.go           # App bootstrap, configuration, and dependency injection
  ├─ routes.go         # HTTP route definitions with middleware chains
  ├─ handlers.go       # HTTP request handlers with form processing
  ├─ middleware.go     # HTTP middleware functions
  ├─ templates.go      # Template functions and cache
  └─ helpers.go        # Shared helper functions and form decoding
internal/models/       # Data models and database operations
  ├─ snippets.go       # Core model definitions with CRUD operations
  └─ errors.go         # Custom error types
internal/validator/    # Input validation framework
  └─ validator.go      # Reusable validation functions and validator struct
ui/                    # User interface components
  ├─ html/             # HTML templates
  │   ├─ base.tmpl     # Base layout template
  │   ├─ pages/        # Page-specific templates (home, view, create)
  │   └─ partials/     # Reusable template components (nav)
  └─ static/           # Static assets
      ├─ css/          # Stylesheets with form styling
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
    formDecoder   *form.Decoder  // Added for form processing
}

// Environment variable handling for sensitive data
password := os.Getenv("DB_PASSWORD")
if password == "" {
    logger.Error("DB_PASSWORD environment variable not set")
    os.Exit(1)
}

// Form decoder initialization
formDecoder := form.NewDecoder()
```

### 2. Form Handling and Validation System

#### Validation Framework Structure

```go
// Validator struct with embedded field errors
type Validator struct {
    FieldErrors map[string]string
}

// Core validation methods
func (v *Validator) Valid() bool {
    return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
    if v.FieldErrors == nil {
        v.FieldErrors = make(map[string]string)
    }
    if _, exists := v.FieldErrors[key]; !exists {
        v.FieldErrors[key] = message
    }
}

func (v *Validator) CheckField(ok bool, key, message string) {
    if !ok {
        v.AddFieldError(key, message)
    }
}

// Reusable validation functions
func NotBlank(value string) bool {
    return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
    return utf8.RuneCountInString(value) <= n
}

func PermittedValues[T comparable](value T, permittedValues ...T) bool {
    return slices.Contains(permittedValues, value)
}
```

#### Form Struct Pattern

```go
// Form struct with embedded validator
type snippetCreateForm struct {
    Title               string `form:"title"`
    Content             string `form:"content"`
    Expires             int    `form:"expires"`
    validator.Validator `form:"-"`
}

// Form processing helper
func (app *application) decodePostForm(r *http.Request, dst any) error {
    err := r.ParseForm()
    if err != nil {
        return err
    }

    err = app.formDecoder.Decode(dst, r.PostForm)
    if err != nil {
        var invalidDecoderError *form.InvalidDecoderError
        if errors.As(err, &invalidDecoderError) {
            panic(err)
        }
        return err
    }
    return nil
}
```

#### Complete Form Handler Pattern

```go
// GET handler - display form
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = snippetCreateForm{
        Expires: 365, // Set default value
    }
    app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// POST handler - process form with validation
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    var form snippetCreateForm

    // Decode form data into struct
    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    // Perform validation
    form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.PermittedValues(form.Expires, 1, 7, 365), "expires", "This field must be one of the following values: 1, 7, or 365")

    // Handle validation errors
    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
        return
    }

    // Process valid form data
    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }
    
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

### 3. Template System with Form Support

```go
// Enhanced template data structure
type templateData struct {
    CurrentYear int
    Snippet     models.Snippet
    Snippets    []models.Snippet
    Form        any  // Generic form field for any form type
}

// Template with form handling and error display
{{define "main"}}
    <form action="/snippet/create" method="post">
        <div>
            <label>Title:</label>
            {{with .Form.FieldErrors.title}}
                <label class="error">{{.}}</label>
            {{end}}
            <input type="text" name="title" value="{{.Form.Title}}">
        </div>
        <div>
            <label>Content:</label>
            {{with .Form.FieldErrors.content}}
                <label class="error">{{.}}</label>
            {{end}}
            <textarea name="content">{{.Form.Content}}</textarea>
        </div>
        <div>
            <label>Delete in:</label>
            {{with .Form.FieldErrors.expires}}
                <label class="error">{{.}}</label>
            {{end}}
            <input type="radio" name="expires" value="365" {{if (eq .Form.Expires 365)}} checked{{end}}> One Year
            <input type="radio" name="expires" value="7" {{if (eq .Form.Expires 7)}} checked {{end}}> One Week
            <input type="radio" name="expires" value="1" {{if (eq .Form.Expires 1)}} checked {{end}}> One Day
        </div>
        <div>
            <input type="submit" value="Create Snippet">
        </div>
    </form>
{{end}}
```

### 4. CSS Styling for Forms

```css
/* Form styling with error states */
form div {
    margin-bottom: 18px;
}

form input[type="text"], form input[type="password"], form input[type="email"] {
    padding: 0.75em 18px;
    width: 100%;
}

form input[type=text], form input[type="password"], form input[type="email"], textarea {
    color: #6A6C6F;
    background: #FFFFFF;
    border: 1px solid #E4E5E7;
    border-radius: 3px;
}

form label {
    display: inline-block;
    margin-bottom: 9px;
}

/* Error styling */
.error {
    color: #C0392B;
    font-weight: bold;
    display: block;
}

.error + textarea, .error + input {
    border-color: #C0392B !important;
    border-width: 2px !important;
}

/* Submit button styling */
input[type="submit"] {
    background-color: #62CB31;
    border-radius: 3px;
    color: #FFFFFF;
    padding: 18px 27px;
    border: none;
    display: inline-block;
    margin-top: 18px;
    font-weight: 700;
}

input[type="submit"]:hover {
    background-color: #4EB722;
    cursor: pointer;
}
```

## Advanced Form Handling Patterns

### 1. Sticky Forms (Value Preservation)

Forms preserve user input when validation fails, improving user experience:

```go
// Template preserves form values
<input type="text" name="title" value="{{.Form.Title}}">
<textarea name="content">{{.Form.Content}}</textarea>
```

### 2. Field-Specific Error Display

Each field can display its own validation errors:

```go
{{with .Form.FieldErrors.title}}
    <label class="error">{{.}}</label>
{{end}}
```

### 3. Conditional Radio Button Selection

Radio buttons maintain selection state based on form data:

```go
<input type="radio" name="expires" value="365" {{if (eq .Form.Expires 365)}} checked{{end}}> One Year
```

### 4. Generic Form Validation

Validation functions use generics for type safety and reusability:

```go
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
    return slices.Contains(permittedValues, value)
}
```

### 5. Professional Error Handling

HTTP status codes properly indicate validation state:

```go
// 422 Unprocessable Entity for validation errors
app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)

// 400 Bad Request for form decoding errors
app.clientError(w, http.StatusBadRequest)
```

## Security Best Practices Implemented

### 1. Server-Side Validation

- All user input is validated on the server
- Client-side validation is never trusted
- Multiple validation rules per field

### 2. Input Sanitization

- UTF-8 aware character counting
- String trimming for blank checks
- Controlled value validation for restricted fields

### 3. Length Limits

- Maximum character limits prevent buffer attacks
- UTF-8 rune counting for accurate character limits
- Database field size alignment

### 4. CSRF Protection Considerations

While not yet implemented, the form structure supports CSRF tokens:

```go
// Future CSRF token field
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
```

## Dependencies and External Libraries

### Core Dependencies

- `github.com/go-sql-driver/mysql` - MySQL database driver
- `github.com/justinas/alice` - HTTP middleware chaining
- `github.com/go-playground/form/v4` - Professional form processing

### Standard Library Usage

- `net/http` - HTTP server and routing
- `html/template` - Template rendering
- `log/slog` - Structured logging
- `database/sql` - Database interface
- `flag` - Command-line argument parsing
- `unicode/utf8` - UTF-8 string processing
- `slices` - Generic slice operations

## Testing Strategies for Forms

### Unit Testing Validation Functions

```go
func TestNotBlank(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"", false},
        {"   ", false},
        {"hello", true},
        {"  hello  ", true},
    }
    
    for _, test := range tests {
        result := validator.NotBlank(test.input)
        if result != test.expected {
            t.Errorf("NotBlank(%q) = %v; want %v", test.input, result, test.expected)
        }
    }
}
```

### Integration Testing Form Handlers

```go
func TestSnippetCreatePost(t *testing.T) {
    app := &application{...} // Initialize test app
    
    form := url.Values{}
    form.Add("title", "Test Title")
    form.Add("content", "Test Content")
    form.Add("expires", "7")
    
    req, _ := http.NewRequest("POST", "/snippet/create", strings.NewReader(form.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    
    rr := httptest.NewRecorder()
    app.snippetCreatePost(rr, req)
    
    if rr.Code != http.StatusSeeOther {
        t.Errorf("Expected status %d, got %d", http.StatusSeeOther, rr.Code)
    }
}
```

## Performance Considerations

### 1. Template Caching

Templates are parsed once at startup for optimal performance

### 2. Form Decoder Reuse

Single form decoder instance shared across requests

### 3. Validation Short-Circuiting

Validation stops at first error per field to reduce processing

### 4. Memory-Efficient Error Storage

Field errors use map for O(1) lookup and minimal memory overhead

## Future Enhancements

### Planned Form Features

- CSRF protection implementation
- File upload handling
- Multi-step forms with session storage
- AJAX form submission with JSON responses
- Client-side validation for better UX

### Advanced Validation

- Email format validation
- Password strength checking
- Custom validation rules
- Cross-field validation (e.g., password confirmation)

### Internationalization

- Multi-language error messages
- Locale-aware validation (date formats, etc.)
- Template translation support

## Common Pitfalls and Solutions

### 1. Form Parsing Errors

**Problem**: Forgetting to call `r.ParseForm()` before accessing form data
**Solution**: Always use the `decodePostForm` helper which handles parsing

### 2. Template Context Issues

**Problem**: Incorrect field access in templates
**Solution**: Use `{{with .Form.FieldErrors.fieldname}}` for conditional error display

### 3. Validation Logic Errors

**Problem**: Client-side validation bypass
**Solution**: Always perform server-side validation regardless of client-side checks

### 4. Memory Leaks in Form Processing

**Problem**: Creating new decoder instances per request
**Solution**: Share single decoder instance across application

## Conclusion

This guide demonstrates a complete, production-ready approach to form handling in Go web applications, featuring:

- Professional validation framework with reusable components
- Clean separation of concerns between validation, processing, and display
- Secure input handling with proper validation and sanitization
- User-friendly error handling with sticky forms
- Comprehensive CSS styling for professional appearance
- Scalable architecture supporting multiple form types

The patterns shown here can be extended to handle complex form requirements while maintaining code quality and security
standards.

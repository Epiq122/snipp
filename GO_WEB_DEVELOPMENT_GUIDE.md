# Go Web Development Project Guide

This document serves as a comprehensive guide for Go web application development, based on the Snippet project. It
outlines key components, best practices, and implementation patterns that can be reused in future projects.

## Project Structure

A well-organized project structure helps maintain code quality and promotes separation of concerns:

```
cmd/web/               # Application entry point and web components
  ├─ main.go           # App bootstrap, configuration, and dependency injection
  ├─ routes.go         # HTTP route definitions with middleware chains and authentication
  ├─ handlers.go       # HTTP request handlers with authentication and form processing
  ├─ middleware.go     # HTTP middleware functions including authentication and CSRF
  ├─ templates.go      # Template functions and cache with authentication support
  └─ helpers.go        # Shared helper functions, form decoding, and authentication helpers
internal/models/       # Data models and database operations
  ├─ snippets.go       # Core snippet model with CRUD operations
  ├─ users.go          # User model with authentication operations
  └─ errors.go         # Custom error types including authentication errors
internal/validator/    # Input validation framework
  └─ validator.go      # Reusable validation functions with authentication support
tls/                   # TLS certificate management
  ├─ cert.pem          # TLS certificate (development)
  └─ key.pem           # TLS private key (development)
ui/                    # User interface components
  ├─ html/             # HTML templates
  │   ├─ base.tmpl     # Base layout template with authentication state
  │   ├─ pages/        # Page-specific templates (home, view, create, signup, login)
  │   └─ partials/     # Reusable template components (nav with authentication)
  └─ static/           # Static assets
      ├─ css/          # Stylesheets with form and authentication styling
      ├─ js/           # JavaScript files
      └─ img/          # Images and icons
```

## Core Components and Implementation Patterns

### 1. Complete User Authentication System

#### User Model Implementation

```go
// User struct representing the data model
type User struct {
ID             int
Name           string
Email          string
HashedPassword []byte
Created        time.Time
}

// UserModel for database operations
type UserModel struct {
DB *sql.DB
}

// User registration with secure password hashing
func (m *UserModel) Insert(name, email, password string) error {
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
if err != nil {
return err
}

stmt := `INSERT INTO users (name, email, hashed_password, created)
             VALUES(?, ?, ?, UTC_TIMESTAMP())`

_, err = m.DB.Exec(stmt, name, email, hashedPassword)
if err != nil {
var mySQLError *mysql.MySQLError
if errors.As(err, &mySQLError) {
if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
return ErrDuplicateEmail
}
}
return err
}
return nil
}

// User authentication with credential verification
func (m *UserModel) Authenticate(email, password string) (int, error) {
var id int
var hashedPassword []byte

stmt := `SELECT id, hashed_password FROM users WHERE email = ?`
err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
if err != nil {
if errors.Is(err, sql.ErrNoRows) {
return 0, ErrInvalidCredentials
}
return 0, err
}

err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
if err != nil {
if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
return 0, ErrInvalidCredentials
}
return 0, err
}

return id, nil
}
```

#### Authentication Error Handling

```go
// Custom authentication errors
var (
ErrNoRecord = errors.New("models: no matching records found")
ErrInvalidCredentials = errors.New("models: invalid credentials provided")
ErrDuplicateEmail = errors.New("models: duplicate email provided")
)
```

### 2. CSRF Protection Implementation

#### CSRF Middleware Configuration

```go
// CSRF protection middleware using nosurf
func preventCSRF(next http.Handler) http.Handler {
csrfHandler := nosurf.New(next)
csrfHandler.SetBaseCookie(http.Cookie{
HttpOnly: true,
Path:     "/",
Secure:   true, // HTTPS-only in production
})
return csrfHandler
}

// Template data structure with CSRF support
type templateData struct {
CurrentYear     int
Snippet         models.Snippet
Snippets        []models.Snippet
Form            any
Flash           string
IsAuthenticated bool
CSRFToken       string // CSRF token for forms
}

// Helper automatically populates CSRF tokens
func (app *application) newTemplateData(r *http.Request) templateData {
return templateData{
CurrentYear:     time.Now().Year(),
Flash:           app.sessionManager.PopString(r.Context(), "flash"),
IsAuthenticated: app.isAuthenticated(r),
CSRFToken:       nosurf.Token(r), // Automatic CSRF token inclusion
}
}
```

#### CSRF Token Usage in Templates

```html
<!-- All forms include CSRF tokens -->
<form action="/user/signup" method="POST" novalidate>
    <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
    <!-- form fields -->
</form>

<form action="/user/logout" method="POST">
    <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
    <button>Logout</button>
</form>
```

### 3. Authentication Middleware and Route Protection

#### Authentication Middleware Implementation

```go
// Authentication middleware for protected routes
func (app *application) requireAuthentication(next http.Handler) http.Handler {
return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
if !app.isAuthenticated(r) {
http.Redirect(w, r, "/user/login", http.StatusSeeOther)
return
}
// Prevent caching of protected content
w.Header().Set("Cache-Control", "no-store")
next.ServeHTTP(w, r)
})
}

// Authentication helper function
func (app *application) isAuthenticated(r *http.Request) bool {
return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
```

#### Multi-layered Route Architecture

```go
// Sophisticated route protection with multiple middleware layers
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

// Static files (no middleware needed)
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// Dynamic routes with session and CSRF protection
dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF)
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

// Public authentication routes
mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

// Protected routes requiring authentication
protected := dynamic.Append(app.requireAuthentication)
mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

// Standard middleware for all routes
standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
return standard.Then(mux)
}
```

### 4. Enhanced Validation Framework with Authentication

#### Comprehensive Validation Functions

```go
// Enhanced validator struct with non-field error support
type Validator struct {
NonFieldErrors []string // For authentication failures
FieldErrors    map[string]string // For field-specific errors
}

// Validation state checking
func (v *Validator) Valid() bool {
return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// Non-field error handling for authentication
func (v *Validator) AddNonFieldError(message string) {
v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// Email validation with regex pattern
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Authentication-specific validation functions
func MinChars(value string, n int) bool {
return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
return rx.MatchString(value)
}

// Validation functions for user authentication
func NotBlank(value string) bool {
return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
return utf8.RuneCountInString(value) <= n
}
```

### 5. Authentication Handler Patterns

#### User Registration Handler

```go
type userSignupForm struct {
Name                string `form:"name"`
Email               string `form:"email"`
Password            string `form:"password"`
validator.Validator `form:"-"`
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
var form userSignupForm

err := app.decodePostForm(r, &form)
if err != nil {
app.clientError(w, http.StatusBadRequest)
return
}

// Comprehensive validation
form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

if !form.Valid() {
data := app.newTemplateData(r)
data.Form = form
app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
return
}

// Handle user creation
err = app.users.Insert(form.Name, form.Email, form.Password)
if err != nil {
if errors.Is(err, models.ErrDuplicateEmail) {
form.AddFieldError("email", "Email address is already in use")
data := app.newTemplateData(r)
data.Form = form
app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
} else {
app.serverError(w, r, err)
}
return
}

app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
```

#### User Login Handler with Session Security

```go
type userLoginForm struct {
Email               string `form:"email"`
Password            string `form:"password"`
validator.Validator `form:"-"`
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
var form userLoginForm

err := app.decodePostForm(r, &form)
if err != nil {
app.clientError(w, http.StatusBadRequest)
return
}

// Validation
form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

if !form.Valid() {
data := app.newTemplateData(r)
data.Form = form
app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
return
}

// Authentication
userID, err := app.users.Authenticate(form.Email, form.Password)
if err != nil {
if errors.Is(err, models.ErrInvalidCredentials) {
form.AddNonFieldError("Email address or password is incorrect")
data := app.newTemplateData(r)
data.Form = form
app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
} else {
app.serverError(w, r, err)
}
return
}

// Session security: renew token to prevent session fixation
err = app.sessionManager.RenewToken(r.Context())
if err != nil {
app.serverError(w, r, err)
return
}

// Store authenticated user ID in session
app.sessionManager.Put(r.Context(), "authenticatedUserID", userID)
http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}
```

#### Secure Logout Handler

```go
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
// Renew session token for security
err := app.sessionManager.RenewToken(r.Context())
if err != nil {
app.serverError(w, r, err)
return
}

// Remove authentication data
app.sessionManager.Remove(r.Context(), "authenticatedUserID")

// Provide user feedback
app.sessionManager.Put(r.Context(), "flash", "You have been logged out successfully.")
http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

### 6. Authentication-Aware Templates

#### Dynamic Navigation Based on Authentication

```html
{{define "nav"}}
<nav>
    <div>
        <a href="/">Home</a>
        {{if .IsAuthenticated}}
        <a href="/snippet/create">Create Snippet</a>
        {{end}}
    </div>
    <div>
        {{if .IsAuthenticated}}
        <form action="/user/logout" method="POST">
            <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
            <button>Logout</button>
        </form>
        {{else}}
        <a href="/user/signup">Signup</a>
        <a href="/user/login">Login</a>
        {{end}}
    </div>
</nav>
{{end}}
```

#### Authentication Forms with Error Handling

```html
<!-- Login form with non-field error support -->
{{define "main"}}
<form action="/user/login" method="POST">
    <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">

    <!-- Display non-field errors (authentication failures) -->
    {{range .Form.NonFieldErrors}}
    <div class="error">{{.}}</div>
    {{end}}

    <div>
        <label>Email:</label>
        {{with .Form.FieldErrors.email}}
        <div class="error">{{.}}</div>
        {{end}}
        <input type="email" name="email" value="{{.Form.Email}}">
    </div>

    <div>
        <label>Password:</label>
        {{with .Form.FieldErrors.password}}
        <div class="error">{{.}}</div>
        {{end}}
        <input type="password" name="password">
    </div>

    <div>
        <input type="submit" value="Login">
    </div>
</form>
{{end}}
```

// ...existing HTTPS/TLS and session management sections continue...

## Authentication Security Best Practices

### 1. Password Security

```go
// Industry-standard password hashing
const bcryptCost = 12
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

// Secure password comparison (constant-time)
err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
```

### 2. Session Security

```go
// Session token renewal on authentication state changes
err = app.sessionManager.RenewToken(r.Context())

// Secure session data storage
app.sessionManager.Put(r.Context(), "authenticatedUserID", userID)

// Session cleanup on logout
app.sessionManager.Remove(r.Context(), "authenticatedUserID")
```

### 3. CSRF Protection

```go
// Automatic CSRF token generation and validation
CSRFToken: nosurf.Token(r)

// Secure CSRF cookie configuration
csrfHandler.SetBaseCookie(http.Cookie{
HttpOnly: true,
Path:     "/",
Secure:   true,
})
```

### 4. Database Security

```go
// Prepared statements prevent SQL injection
stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

// Proper error handling for database constraints
if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
return ErrDuplicateEmail
}
```

## Testing Authentication Systems

### 1. Unit Testing Authentication Functions

```go
func TestUserAuthenticate(t *testing.T) {
// Test valid credentials
userID, err := userModel.Authenticate("valid@example.com", "validpassword")
if err != nil {
t.Errorf("Expected successful authentication, got error: %v", err)
}

// Test invalid credentials
_, err = userModel.Authenticate("invalid@example.com", "wrongpassword")
if !errors.Is(err, models.ErrInvalidCredentials) {
t.Errorf("Expected ErrInvalidCredentials, got: %v", err)
}
}
```

### 2. Integration Testing Authentication Flows

```go
func TestUserSignupFlow(t *testing.T) {
form := url.Values{}
form.Add("name", "Test User")
form.Add("email", "test@example.com")
form.Add("password", "testpassword")
form.Add("csrf_token", "valid-token")

req := httptest.NewRequest("POST", "/user/signup", strings.NewReader(form.Encode()))
req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

rr := httptest.NewRecorder()
app.userSignupPost(rr, req)

if rr.Code != http.StatusSeeOther {
t.Errorf("Expected redirect after signup, got %d", rr.Code)
}
}
```

## Dependencies and External Libraries

### Authentication Dependencies

- `golang.org/x/crypto/bcrypt` - Secure password hashing
- `github.com/justinas/nosurf` - CSRF protection middleware
- `github.com/go-sql-driver/mysql` - Database operations
- `github.com/alexedwards/scs/v2` - Session management

### Integration Benefits

- Secure password storage with industry-standard hashing
- Complete CSRF protection for all state-changing operations
- Professional session management with database storage
- Comprehensive validation framework with authentication support

## Conclusion

The authentication system represents the final major milestone in creating a production-ready web application with Go.
The implementation provides:

- **Complete User Management**: Registration, login, logout with secure password handling
- **Comprehensive Security**: CSRF protection, session security, input validation
- **Professional UX**: Dynamic navigation, form validation, user feedback
- **Scalable Architecture**: Multi-layered middleware, clean separation of concerns
- **Production-Ready**: Industry-standard security practices and error handling

This foundation enables building complex user-centric web applications while maintaining the highest security standards
and providing an excellent user experience.

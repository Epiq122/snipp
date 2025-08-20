package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define errors that the model might return
var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert adds a new snippet to the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Using standard ANSI SQL for better compatibility
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Get retrieves a specific snippet by ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// SQL statement to select snippet
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE id = ? AND expires > UTC_TIMESTAMP()`

	// Use QueryRow() to execute the SQL statement, passing in the id value
	// as the placeholder parameter
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// Use row.Scan() to copy the values from each field in the row to the
	// corresponding field in the Snippet struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, row.Scan() will return sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	// If everything went OK, return the Snippet object
	return s, nil
}

// Latest retrieves the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// SQL statement to retrieve the 10 most recent snippets
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use Query() to execute the SQL statement
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize an empty slice to hold the snippets
	snippets := []*Snippet{}

	// Iterate through the rows in the result set
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct
		s := &Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to
		// the corresponding field in the Snippet struct
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append the Snippet to our slice
		snippets = append(snippets, s)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

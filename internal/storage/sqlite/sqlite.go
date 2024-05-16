// Package sqlite provides a SQLite implementation of the storage interface.
package sqlite

import (
	"citizen_webservice/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // Importing the SQLite driver
)

// Storage struct represents a SQLite database.
type Storage struct {
	db *sql.DB
}

// New function initializes a new SQLite database at the provided storage path.
// It returns a pointer to a Storage struct or an error.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	// Open a new database connection
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Check the database connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Prepare a SQL statement to create the users table if it doesn't exist
	stmt, err := db.Prepare(`
 CREATE TABLE IF NOT EXISTS users (
  iin VARCHAR(14) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  phone VARCHAR(30) NOT NULL UNIQUE
 );`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Execute the SQL statement
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Return a new Storage struct
	return &Storage{db: db}, nil
}

// SavePerson method saves a person's information in the database.
// It returns an error if the operation fails.
func (s *Storage) SavePerson(iin string, name string, phone string) error {
	const op = "storage.sqlite.SavePerson"

	// Prepare a SQL statement to insert a new user
	stmt, err := s.db.Prepare("INSERT INTO users(iin, name, phone) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Execute the SQL statement
	_, err = stmt.Exec(iin, name, phone)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
				return fmt.Errorf("%s: %w", op, storage.ErrorIINExists)
			}
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return fmt.Errorf("%s: %w", op, storage.ErrorPhoneNumberExists)
			}

		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetPersonByIIN method retrieves a person's information by their IIN.
// It returns a PersonInfo struct or an error.
func (s *Storage) GetPersonByIIN(iin string) (storage.PersonInfo, error) {
	const fn = "storage.sqlite.GetPersonByIIN"

	// Prepare a SQL statement to select a user by IIN
	stmt, err := s.db.Prepare("SELECT * FROM users WHERE iin = ? LIMIT 1;")
	if err != nil {
		return storage.PersonInfo{}, fmt.Errorf("%s: %w", fn, err)
	}
	person := storage.PersonInfo{}
	err = stmt.QueryRow(iin).Scan(&person.IIN, &person.Name, &person.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.PersonInfo{}, fmt.Errorf("%s: %w", fn, storage.ErrorIINNotFound)
		}
		return storage.PersonInfo{}, fmt.Errorf("%s: %w", fn, err)
	}
	return person, nil
}

// GetPersonByName method retrieves all people with a name that matches the provided name.
// It returns a slice of PersonInfo structs or an error.
func (s *Storage) GetPersonByName(name string) ([]storage.PersonInfo, error) {
	const fn = "storage.sqlite.GetPersonByName"
	var allMatchedPeople []storage.PersonInfo

	// Prepare a SQL statement to select users by name
	stmt, err := s.db.Prepare("SELECT * FROM users WHERE name LIKE ?;")
	if err != nil {
		return allMatchedPeople, fmt.Errorf("%s: %w", fn, err)
	}

	// Execute the SQL statement
	rows, err := stmt.Query("%" + name + "%")
	if err != nil {
		return allMatchedPeople, fmt.Errorf("%s: %w", fn, err)
	}

	// Scan the result rows into PersonInfo structs
	for rows.Next() {
		person := storage.PersonInfo{}
		err = rows.Scan(&person.IIN, &person.Name, &person.Phone)
		if err != nil {
			return allMatchedPeople, fmt.Errorf("%s: %w", fn, err)
		}
		allMatchedPeople = append(allMatchedPeople, person)
	}

	if err = rows.Err(); err != nil {
		return allMatchedPeople, fmt.Errorf("%s: %w", fn, err)
	}

	return allMatchedPeople, nil
}

// DeletePersonByIIN method deletes a person's information by their IIN.
// It returns an error if the operation fails.
func (s *Storage) DeletePersonByIIN(iin string) error {
	const fn = "storage.sqlite.DeletePersonByIIN"

	// Prepare a SQL statement to delete a user by IIN
	stmt, err := s.db.Prepare("DELETE FROM users WHERE iin = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	// Execute the SQL statement
	result, err := stmt.Exec(iin)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", fn, storage.ErrorIINNotFound)
	}

	return nil
}

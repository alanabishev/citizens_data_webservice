package sqlite

import (
	"citizen_webservice/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		iin VARCHAR(14) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		phone VARCHAR(30) NOT NULL UNIQUE
	);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SavePerson(iin string, name string, phone string) error {
	const op = "storage.sqlite.SavePerson"

	stmt, err := s.db.Prepare("INSERT INTO users(iin, name, phone) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(iin, name, phone)
	if err != nil {
		//if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		//	return fmt.Errorf("%s: %w", op, storage.ErrorIINExists)
		//}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetPersonByIIN(iin string) (storage.PersonInfo, error) {
	const fn = "storage.sqlite.GetPersonByIIN"

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

func (s *Storage) GetPersonByName(name string) ([]storage.PersonInfo, error) {
	const fn = "storage.sqlite.GetPersonByName"
	var allMatchedPeople []storage.PersonInfo

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE name LIKE ?;")
	if err != nil {
		return allMatchedPeople, fmt.Errorf("%s: %w", fn, err)
	}

	rows, err := stmt.Query("%" + name + "%")
	if err != nil {
		return allMatchedPeople, fmt.Errorf("%s: %w", fn, err)
	}

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

func (s *Storage) DeletePersonByIIN(iin string) error {
	const op = "storage.sqlite.DeletePersonByIIN"

	stmt, err := s.db.Prepare("DELETE FROM users WHERE iin = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(iin)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

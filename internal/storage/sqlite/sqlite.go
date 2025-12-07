// package sqlite

// import (
// 	"database/sql"

// 	config "github.com/akshayjha21/Student-Api/internal/config"
// 	_ "github.com/mattn/go-sqlite3"
// )

// type Sqlite struct {
// 	Db *sql.DB
// }

// func New(cfg *config.Config) (*Sqlite, error) {

//		db, err := sql.Open("sqlite3",config.Config.StoragePath)
//		if err != nil {
//			return nil, err
//		}
//		_, err = db.Exec(`CREATE TABLE IF NOT EXIST students(
//		id INTEGER PRIMARY KEY AUTOINCREMENT
//		name TEXT,
//		email TEXT,
//		age INTEGER
//		)`)
//		if err != nil {
//			return nil, err
//		}
//		return &Sqlite{
//			Db: db,
//		}, nil
//	}
package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	config "github.com/akshayjha21/Student-Api/internal/config"
	"github.com/akshayjha21/Student-Api/internal/types"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			age INTEGER
		);
	`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) CreateStudent(name, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare(`
		INSERT INTO students (name, email, age)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

}

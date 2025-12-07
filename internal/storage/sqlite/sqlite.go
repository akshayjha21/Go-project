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

	// "golang.org/x/tools/go/analysis/passes/defers"
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

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) UpdateById(id int64, data types.Student) (types.Student, error) {
	stmt, err := s.Db.Prepare(`
		UPDATE students
		SET name=?,email=?,age=?
		WHERE id=? 
	`)
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(data.Name, data.Email, data.Age,id)
	if err != nil {
		return types.Student{}, fmt.Errorf("update error: %w", err)
	}
rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.Student{}, err
	}

	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id %d", id)
	}

	// return updated student
	return s.GetStudentById(id)
}

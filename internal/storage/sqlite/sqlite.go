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

	config "github.com/akshayjha21/Student-Api/internal/config"
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

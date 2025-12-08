package storage

import (
	"github.com/akshayjha21/Student-Api/internal/types"
	"github.com/akshayjha21/Student-Api/internal/utils/pagination"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents(p*pagination.Paginate) ([]types.Student, error)
	UpdateById(id int64, data types.Student) (types.Student, error)
	DeleteByID(id int64) error
	UpdateField(id int64, data types.StudentPatch) (types.Student, error)
}

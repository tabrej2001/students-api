package storage

import "github.com/tabrej2001/student-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Students, error)
}

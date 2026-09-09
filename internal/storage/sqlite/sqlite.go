package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tabrej2001/student-api/internal/config"
	"github.com/tabrej2001/student-api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.Storage_path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, nil
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Students, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age from Students WHERE id = ? lIMIT 1")
	if err != nil {
		return types.Students{}, nil
	}

	defer stmt.Close()

	var student types.Students

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Students{}, fmt.Errorf("No student found with id %s", fmt.Sprint(id))
		}

		return types.Students{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudentList() ([]types.Students, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age from Students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Students

	for rows.Next() {
		var student types.Students
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) UpdateStudent(id int64, student types.Students) (types.Students, error) {
	stmt, err := s.Db.Prepare(`
	UPDATE students 
	SET name = ?, email = ?, age = ?
	WHERE id = ?
	`)
	if err != nil {
		return types.Students{}, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		student.Name,
		student.Email,
		student.Age,
		id,
	)
	if err != nil {
		return types.Students{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.Students{}, err
	}

	if rowsAffected == 0 {
		return types.Students{}, fmt.Errorf("no student with id %d", id)
	}

	updatedStudent, err := s.GetStudentById(id)
	if err != nil {
		return types.Students{}, err
	}

	return updatedStudent, nil
}

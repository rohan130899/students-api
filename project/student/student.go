package student

import (
	"context"
	"errors"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrNoStudentFound  = errors.New("no student found")
	ErrUpdatingStudent = errors.New("failed to update student")
	ErrDeletingStudent = errors.New("failed to delete student")
	ErrFetchingStudent = errors.New("failed to fetch student by ID")
)

type CustomeTime struct {
	time.Time
}

const dateLayout = "02-01-2006"

func (ct *CustomeTime) UnmarshalJSON(b []byte) (err error) {
	dateStr := strings.Trim(string(b), `"`)
	ct.Time, err = time.Parse(dateLayout, dateStr)
	return
}

type Student struct {
	ID        string
	Name      string
	Email     string
	Address   string
	Gender    string
	CreatedBy string
	CreatedOn time.Time
	UpdatedBy string
	UpdatedOn time.Time
}

type StudentStore interface {
	CreateStudent(context.Context, Student) (Student, error)
	GetStudent(context.Context, string) (Student, error)
	DeleteStudent(context.Context, string) error
	UpdateStudent(context.Context, string, Student) (Student, error)
	GetStudents(context.Context) ([]Student, error)
	Ping(context.Context) error
}

type Service struct {
	db StudentStore
}

func NewService(db StudentStore) *Service {
	return &Service{
		db: db,
	}
}

// CreateStudent: Creates a Student
func (s *Service) CreateStudent(ctx context.Context, student Student) (Student, error) {
	student, err := s.db.CreateStudent(ctx, student)
	if err != nil {
		log.Errorf("an error occurred adding the Student: %s", err.Error())
	}
	return student, nil

}

// GetStudent: Retrieves Student by their ID
func (s *Service) GetStudent(ctx context.Context, ID string) (Student, error) {
	// calls store passing in the context
	std, err := s.db.GetStudent(ctx, ID)
	if err != nil {
		log.Errorf("an error occured fetching the Student: %s", err.Error())
		return Student{}, ErrFetchingStudent
	}
	return std, nil
}

// DeleteStudent: Deletes a student by ID
func (s *Service) DeleteStudent(ctx context.Context, ID string) error {
	return s.db.DeleteStudent(ctx, ID)
}

// UpdateStudent: Updates a student by ID with new info
func (s *Service) UpdateStudent(ctx context.Context, ID string, newStudent Student) (Student, error) {
	std, err := s.db.UpdateStudent(ctx, ID, newStudent)
	if err != nil {
		log.Errorf("an error occurred updating the Student: %s", err.Error())
	}
	return std, nil
}

// GetStudents: Get all student details
func (s *Service) GetStudents(ctx context.Context) ([]Student, error) {
	students, err := s.db.GetStudents(ctx)
	if err != nil {
		log.Errorf("an error occurred updating the Student: %s", err.Error())
	}
	return students, nil
}

func (s *Service) ReadyCheck(ctx context.Context) error {
	log.Info("Checking Readiness")
	return s.db.Ping(ctx)
}

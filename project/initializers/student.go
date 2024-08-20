package initializers

import (
	"context"
	"database/sql"
	"fmt"
	"students-api/project/authentication"
	"students-api/project/student"
	"time"

	log "github.com/sirupsen/logrus"

	uuid "github.com/satori/go.uuid"
)

type StudentRow struct {
	ID        string
	Name      sql.NullString
	Email     sql.NullString
	Gender    sql.NullString
	Address   sql.NullString
	CreatedBy sql.NullString
	CreatedOn sql.NullTime
	UpdatedBy sql.NullString
	UpdatedOn sql.NullTime
}

func convertStudentRowToStudent(s StudentRow) student.Student {
	return student.Student{
		ID:        s.ID,
		Name:      s.Name.String,
		Email:     s.Email.String,
		Gender:    s.Gender.String,
		Address:   s.Address.String,
		CreatedBy: s.CreatedBy.String,
		CreatedOn: s.CreatedOn.Time,
		UpdatedBy: s.UpdatedBy.String,
		UpdatedOn: s.UpdatedOn.Time,
	}
}

func (d *Database) CreateStudent(ctx context.Context, std student.Student) (student.Student, error) {
	createdby, ok := ctx.Value(authentication.UserIDKey).(string)
	log.Info(createdby)
	if !ok {
		return student.Student{}, fmt.Errorf("could not retrieve user ID from context")
	}
	std.ID = uuid.NewV4().String()
	std.CreatedBy = createdby
	postRow := StudentRow{
		ID:        std.ID,
		Name:      sql.NullString{String: std.Name, Valid: true},
		Email:     sql.NullString{String: std.Email, Valid: true},
		Gender:    sql.NullString{String: std.Gender, Valid: true},
		Address:   sql.NullString{String: std.Address, Valid: true},
		CreatedBy: sql.NullString{String: std.CreatedBy, Valid: true},
		CreatedOn: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy: sql.NullString{String: std.UpdatedBy, Valid: true},
		UpdatedOn: sql.NullTime{Time: time.Now(), Valid: false},
	}

	rows, err := d.Client.NamedExecContext(
		ctx,
		`INSERT INTO student 
 (id, name, email, gender, address, createdby, createdon, updatedby, updatedon) 
 VALUES
 (:id, :name, :email, :gender, :address, :createdby, :createdon, :updatedby, :updatedon)`,
		postRow,
	)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to insert Student: %w", err)
	}
	insertedID, err := rows.LastInsertId()
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to close rows: %w", err)
	}
	fmt.Printf("New record ID: %d\n", insertedID)
	return std, nil
}

func (d *Database) GetStudent(ctx context.Context, uuid string) (student.Student, error) {
	// fetching studentRow from the database and then converting to student.Student
	var stdRow StudentRow
	row := d.Client.QueryRowContext(
		ctx,
		`SELECT id, name, email, gender, address, createdby, createdon, updatedby, updatedon
 FROM student
 WHERE id = ?`,
		uuid,
	)
	err := row.Scan(&stdRow.ID, &stdRow.Name, &stdRow.Email, &stdRow.Gender, &stdRow.Address, &stdRow.CreatedBy, &stdRow.CreatedOn, &stdRow.UpdatedBy, &stdRow.UpdatedOn)
	if err != nil {
		return student.Student{}, fmt.Errorf("an error occurred fetching the student by uuid: %w", err)
	}
	// sqlx with context to ensure context cancelation is honoured
	return convertStudentRowToStudent(stdRow), nil
}

func (d *Database) GetStudents(ctx context.Context) ([]student.Student, error) {
	// Fetching all rows from the database
	var stdRows []StudentRow
	err := d.Client.Select(&stdRows, "SELECT * FROM student LIMIT 10")
	if err != nil {
		return nil, fmt.Errorf("fetchStudents %v", err)
	}

	// Converting each StudentRow to student.Student
	students := make([]student.Student, len(stdRows))
	for i, stdRow := range stdRows {
		students[i] = student.Student{
			ID:        stdRow.ID,
			Name:      stdRow.Name.String,
			Email:     stdRow.Email.String,
			Gender:    stdRow.Gender.String,
			Address:   stdRow.Address.String,
			CreatedBy: stdRow.CreatedBy.String,
			CreatedOn: stdRow.CreatedOn.Time,
			UpdatedBy: stdRow.UpdatedBy.String,
			UpdatedOn: stdRow.UpdatedOn.Time,
		}
	}

	return students, nil
}

func (d *Database) DeleteStudent(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM student WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete comment from the database: %w", err)
	}
	return nil
}

func (d *Database) UpdateStudent(ctx context.Context, id string, std student.Student) (student.Student, error) {
	updatedby, ok := ctx.Value(authentication.UserIDKey).(string)
	std.UpdatedBy = updatedby
	if !ok {
		return student.Student{}, fmt.Errorf("could not retrieve user ID from context")
	}
	stdRow := StudentRow{
		ID:        id,
		Name:      sql.NullString{String: std.Name, Valid: true},
		Email:     sql.NullString{String: std.Email, Valid: true},
		Gender:    sql.NullString{String: std.Gender, Valid: true},
		Address:   sql.NullString{String: std.Address, Valid: true},
		CreatedBy: sql.NullString{String: std.CreatedBy, Valid: true},
		CreatedOn: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy: sql.NullString{String: std.UpdatedBy, Valid: true},
		UpdatedOn: sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows, err := d.Client.NamedExecContext(
		ctx,
		`UPDATE student 
 SET Name = :Name, email = :email, gender = :gender, address = :address, createdby = :createdby, createdon = :createdon, updatedby = :updatedby, updatedon = :updatedon WHERE id = :id`,
		stdRow,
	)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to insert Student: %w", err)
	}

	insertedID, err := rows.LastInsertId()
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to close rows: %w", err)
	}
	fmt.Printf("New record ID: %d\n", insertedID)

	return convertStudentRowToStudent(stdRow), nil
}

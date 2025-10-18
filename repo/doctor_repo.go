package repo

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrDoctorNotFound = errors.New("doctor not found")
	ErrFailedToUpdate = errors.New("failed to update doctor")
	ErrFailedToDelete = errors.New("failed to delete doctor")
)

type Doctor struct {
	DoctorID        int       `json:"doctor_id" db:"doctor_id"`
	Name            string    `json:"name" db:"name"`
	Specialty       string    `json:"specialty" db:"specialty"`
	YearsExperience int       `json:"years_experience" db:"years_experience"`
	PhoneNumber     string    `json:"phone_number" db:"phone_number"`
	Email           string    `json:"email" db:"email"`
	ImageURL        string    `json:"image_url" db:"image_url"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type DoctorRepo interface {
	Create(doctor Doctor) (*Doctor, error)
	List(search string, offset, limit int) ([]Doctor, int, error)
	Get(id int) (*Doctor, error)
	Update(doctor Doctor) (*Doctor, error)
	Delete(id int) error
}

type doctorRepo struct {
	db *sqlx.DB
}

func NewDoctorRepo(db *sqlx.DB) DoctorRepo {
	return &doctorRepo{db: db}
}

func (r *doctorRepo) Create(d Doctor) (*Doctor, error) {
	query := `
		INSERT INTO doctors (
		  name,
		  specialty,
		  years_experience,
		  phone_number,
		  email,
		  image_url
		) VALUES (
		   :name,
		   :specialty,
		   :years_experience,
		   :phone_number,
		   :email,
		   :image_url
		)
		RETURNING
		  doctor_id,
		  name,
		  specialty,
		  years_experience,
		  phone_number, email,
		  image_url, created_at,
		  updated_at;
	`
	rows, err := r.db.NamedQuery(query, d)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var created Doctor
		if err := rows.StructScan(&created); err != nil {
			return nil, err
		}
		return &created, nil
	}
	return nil, nil
}

func (r *doctorRepo) List(search string, offset, limit int) ([]Doctor, int, error) {
	var doctors []Doctor
	// search pattern
	searchQuery := "%"
	if search != "" {
		searchQuery = "%" + search + "%"
	}
	var total int
	errCount := r.db.Get(&total, `SELECT COUNT(*) FROM doctors WHERE name ILIKE $1`, searchQuery)
	if errCount != nil {
		return nil, 0, fmt.Errorf("error counting hospitals: %w", errCount)
	}
	query := `
	  SELECT *
	  FROM doctors
	  WHERE name ILIKE $1
	  ORDER BY created_at DESC
	  LIMIT $2 OFFSET $3 
	`
	err := r.db.Select(&doctors, query, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching doctors: %w", err)
	}
	return doctors, total, nil
}

func (r *doctorRepo) Get(id int) (*Doctor, error) {
	var doctor Doctor
	query := `SELECT * FROM doctors WHERE doctor_id = $1`
	err := r.db.Get(&doctor, query, id)
	if err != nil {
		return nil, ErrDoctorNotFound
	}
	return &doctor, nil
}

func (r *doctorRepo) Update(d Doctor) (*Doctor, error) {
	query := `
		UPDATE doctors
		SET 
		  name = :name,
		  specialty = :specialty,
		  years_experience = :years_experience,
		  phone_number = :phone_number,
		  email = :email,
		  image_url = :image_url,
		  updated_at = NOW()
		WHERE doctor_id = :doctor_id
		RETURNING
		  doctor_id,
		  name,
		  specialty,
		  years_experience,
		  phone_number,
		  email,
		  image_url,
		  created_at,
		  updated_at;
	`
	rows, err := r.db.NamedQuery(query, d)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var updated Doctor
		if err := rows.StructScan(&updated); err != nil {
			return nil, err
		}
		return &updated, nil
	}
	return nil, ErrFailedToUpdate
}

func (r *doctorRepo) Delete(id int) error {
	query := `DELETE FROM doctors WHERE doctor_id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrFailedToDelete
	}
	return nil
}

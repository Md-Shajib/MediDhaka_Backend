package repo

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type HospitalDoctor struct {
	HospitalID int       `json:"hospital_id" db:"hospital_id"`
	DoctorID   int       `json:"doctor_id" db:"doctor_id"`
	Role       string    `json:"role" db:"role"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type HospitalDoctorRepo interface {
	AssignDoctor(rel HospitalDoctor) error
	ListDoctorsByHospital(hospitalID int) ([]Doctor, error)
}

type hospitalDoctorRepo struct {
	db *sqlx.DB
}

func NewHospitalDoctorRepo(db *sqlx.DB) HospitalDoctorRepo {
	return &hospitalDoctorRepo{db: db}
}

func (r *hospitalDoctorRepo) AssignDoctor(rel HospitalDoctor) error {
	query := `
		INSERT INTO hospital_doctor (
		  hospital_id,
		  doctor_id,
		  role
		) VALUES (
		  :hospital_id,
		  :doctor_id,
		  :role
		)
	`
	_, err := r.db.NamedExec(query, rel)
	return err
}

func (r *hospitalDoctorRepo) ListDoctorsByHospital(hospitalID int) ([]Doctor, error) {
	var doctors []Doctor
	query := `
		SELECT d.*
		FROM doctors d
		JOIN hospital_doctor hd ON d.doctor_id = hd.doctor_id
		WHERE hd.hospital_id = $1
	`
	err := r.db.Select(&doctors, query, hospitalID)
	return doctors, err
}

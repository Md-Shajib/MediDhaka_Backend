package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// Custom public errors
var (
	ErrNotFound     = errors.New("record not found in the database")
	ErrFailedUpdate = errors.New("failed to update record: zero rows affected")
)

// DB structure for a hospital record.
type Hospital struct {
	HospitalID  int       `json:"hospital_id" db:"hospital_id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// HospitalRepo interface.
type hospitalRepo struct {
	dbCon *sqlx.DB
}

// HospitalRepo defines the interface for data access operations.
type HospitalRepo interface {
	Create(hospital Hospital) (*Hospital, error)
	Get(id int) (*Hospital, error)
	List(search string, offset, limit int) ([]*Hospital, int, error)
	Update(h Hospital) (*Hospital, error)
	Delete(id int) error
}

// NewHospitalRepo creates a new repository instance.
func NewHospitalRepo(dbCon *sqlx.DB) HospitalRepo {
	return &hospitalRepo{
		dbCon: dbCon,
	}
}

// INSERT query
func (r *hospitalRepo) Create(hospital Hospital) (*Hospital, error) {
	query := `
		INSERT INTO hospitals (
			name, 
			address, 
			phone_number, 
			email,
			image_url
		)
		VALUES (
			:name, 
			:address, 
			:phone_number, 
			:email,
			:image_url
		)
		RETURNING
		  hospital_id,
		  name, address,
		  phone_number,
		  email,
		  image_url,
		  created_at,
		  updated_at;
	`
	rows, err := r.dbCon.NamedQuery(query, hospital)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdHospital Hospital // Use a new variable name to avoid shadowing
		if err := rows.StructScan(&createdHospital); err != nil {
			return nil, err
		}
		return &createdHospital, nil
	}

	return nil, errors.New("failed to return created hospital data")
}

// Get a single Hospital record by ID.
func (r *hospitalRepo) Get(id int) (*Hospital, error) {
	var hsp Hospital
	query := `SELECT * FROM hospitals WHERE hospital_id = $1`
	err := r.dbCon.Get(&hsp, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound // Return custom public error
		}
		return nil, fmt.Errorf("error fetching hospital: %w", err)
	}

	return &hsp, nil
}

// GET all Hospital records.
func (r *hospitalRepo) List(search string, offset, limit int) ([]*Hospital, int, error) {
	var hspList []*Hospital
	// search pattern
	searchQuery := "%"
	if search != "" {
		searchQuery = "%" + search + "%"
	}

	var total int
	err := r.dbCon.Get(&total, "SELECT COUNT(*) FROM hospitals WHERE name ILIKE $1", searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting hospitals: %w", err)
	}

	query := `
		SELECT *
		FROM hospitals
		WHERE name ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err = r.dbCon.Select(&hspList, query, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching hospitals: %w", err)
	}
	return hspList, total, nil
}

// Update an existing Hospital record.
func (r *hospitalRepo) Update(h Hospital) (*Hospital, error) {
	h.UpdatedAt = time.Now()

	query := `
		UPDATE hospitals
		SET 
		  name = :name,
		  address = :address,
		  phone_number = :phone_number,
		  email = :email,
		  image_url = :image_url,
		  updated_at = :updated_at
		WHERE hospital_id = :hospital_id
		RETURNING 
		  hospital_id,
		  name,
		  address,
		  phone_number,
		  email,
		  image_url,
		  created_at,
		  updated_at;
	`
	rows, err := r.dbCon.NamedQuery(query, h)
	if err != nil {
		return nil, fmt.Errorf("error executing update query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var updatedHospital Hospital
		if err := rows.StructScan(&updatedHospital); err != nil {
			return nil, fmt.Errorf("error scanning updated hospital: %w", err)
		}
		return &updatedHospital, nil
	}

	// If rows.Next() is false, it means no rows were returned, so the ID didn't match.
	return nil, ErrFailedUpdate
}

// Delete a Hospital record by ID.
func (r *hospitalRepo) Delete(id int) error {
	query := `DELETE from hospitals WHERE hospital_id = $1`
	result, err := r.dbCon.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error executing delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Could not get rows affected after delete: %v", err)
		return nil
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

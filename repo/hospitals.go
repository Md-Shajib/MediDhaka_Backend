// package repo

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"time"

// 	"github.com/jmoiron/sqlx"
// )

// type Hospital struct {
// 	HospitalID  int       `json:"hospital_id" db:"hospital_id"`
// 	Name        string    `json:"name" db:"name"`
// 	Address     string    `json:"address" db:"address"`
// 	PhoneNumber string    `json:"phone_number" db:"phone_number"`
// 	Email       string    `json:"email" db:"email"`
// 	CreatedAt   time.Time `json:"created_at" db:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
// }

// type hospitalRepo struct {
// 	dbCon *sqlx.DB
// }

// type HospitalRepo interface {
// 	Create(hospital Hospital) (*Hospital, error)
// 	Get(id int) (*Hospital, error)
// 	List() ([]*Hospital, error)
// 	Update(h Hospital) (*Hospital, error)
// 	Delete(id int) error
// }

// func NewHospitalRepo(dbCon *sqlx.DB) HospitalRepo {
// 	return &hospitalRepo{
// 		dbCon: dbCon,
// 	}
// }

// func (r *hospitalRepo) Create(hospital Hospital) (*Hospital, error) {
// 	query := `
// 		INSERT INTO hospitals (
// 			name,
// 			address,
// 			phone_number,
// 			email
// 		)
// 		VALUES (
// 			:name,
// 			:address,
// 			:phone_number,
// 			:email
// 		)
// 		RETURNING hospital_id, name, address, phone_number, email, created_at, updated_at
// 	`

// 	// sqlx.NamedQuery helps bind struct fields directly using `db:"field"`
// 	rows, err := r.dbCon.NamedQuery(query, hospital)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	if rows.Next() {
// 		var hospital Hospital
// 		if err := rows.StructScan(&hospital); err != nil {
// 			return nil, err
// 		}
// 		fmt.Println("Created Successfully: ", hospital)
// 		return &hospital, nil
// 	}

// 	return nil, nil
// }

// func (r *hospitalRepo) Get(id int) (*Hospital, error) {
// 	var hsp Hospital
// 	query := `
// 	 SELECT
// 	   hospital_id,
// 	   name, address,
// 	   phone_number,
// 	   email,
// 	   created_at,
// 	   updated_at
// 	 FROM hospitals
// 	 WHERE hospital_id = $1
// 	`
// 	err := r.dbCon.Get(&hsp, query, id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("No data available in database")
// 		}
// 		return nil, err
// 	}

// 	return &hsp, nil
// }

// func (r *hospitalRepo) List() ([]*Hospital, error) {
// 	var hspList []*Hospital
// 	query := `
// 	  SELECT
// 	    hospital_id,
// 	    name, address,
// 	    phone_number,
// 	    email,
// 	    created_at,
// 	    updated_at
// 	  FROM hospitals
// 	`
// 	err := r.dbCon.Select(&hspList, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return hspList, nil
// }

// func (r *hospitalRepo) Update(h Hospital) (*Hospital, error) {
// 	query := `
// 	  UPDATE hospitals
// 	  SET
// 		name = :name,
// 		address = :address,
// 		phone_number = :phone_number,
// 		email = :email,
// 		updated_at = NOW() -- Set updated_at automatically
// 	  WHERE hospital_id = :hospital_id
// 	  RETURNING hospital_id, name, address, phone_number, email, created_at, updated_at
// 	`

// 	row := r.dbCon.QueryRow(query, h.Name, h.Address, h.PhoneNumber, h.Email, h.UpdatedAt)
// 	err := row.Err()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &h, nil
// }

// func (r *hospitalRepo) Delete(id int) error {
// 	query := `DELETE from hospitals WHERE hospital_id = $1`

// 	_, err := r.dbCon.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// Custom public errors for the repository layer
var (
	ErrNotFound     = errors.New("record not found in the database")
	ErrFailedUpdate = errors.New("failed to update record: zero rows affected")
)

// Hospital represents the database structure for a hospital record.
type Hospital struct {
	HospitalID  int       `json:"hospital_id" db:"hospital_id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// hospitalRepo implements the HospitalRepo interface.
type hospitalRepo struct {
	dbCon *sqlx.DB
}

// HospitalRepo defines the interface for data access operations.
type HospitalRepo interface {
	Create(hospital Hospital) (*Hospital, error)
	Get(id int) (*Hospital, error)
	List() ([]*Hospital, error)
	Update(h Hospital) (*Hospital, error)
	Delete(id int) error
}

// NewHospitalRepo creates a new repository instance.
func NewHospitalRepo(dbCon *sqlx.DB) HospitalRepo {
	return &hospitalRepo{
		dbCon: dbCon,
	}
}

// Create executes an INSERT query and returns the created record, including auto-generated fields.
func (r *hospitalRepo) Create(hospital Hospital) (*Hospital, error) {
	query := `
		INSERT INTO hospitals (
			name, 
			address, 
			phone_number, 
			email
		)
		VALUES (
			:name, 
			:address, 
			:phone_number, 
			:email
		)
		RETURNING hospital_id, name, address, phone_number, email, created_at, updated_at
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
	// This path should ideally not be reached if RETURNING is used correctly,
	// but return an unexpected error just in case.
	return nil, errors.New("failed to return created hospital data")
}

// Get retrieves a single Hospital record by its ID.
func (r *hospitalRepo) Get(id int) (*Hospital, error) {
	var hsp Hospital
	query := `
	 SELECT 
	 	hospital_id, 
	 	name, address, 
	 	phone_number, 
	 	email, 
	 	created_at, 
	 	updated_at 
	 FROM hospitals 
	 WHERE hospital_id = $1
	`
	err := r.dbCon.Get(&hsp, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound // Return custom public error
		}
		return nil, fmt.Errorf("error fetching hospital: %w", err)
	}

	return &hsp, nil
}

// List retrieves all Hospital records.
func (r *hospitalRepo) List() ([]*Hospital, error) {
	var hspList []*Hospital
	query := `
	 SELECT 
	 	hospital_id, 
	 	name, address, 
	 	phone_number, 
	 	email, 
	 	created_at, 
	 	updated_at 
	 FROM hospitals 
	`
	err := r.dbCon.Select(&hspList, query)
	if err != nil {
		return nil, fmt.Errorf("error listing hospitals: %w", err)
	}
	return hspList, nil
}

// Update modifies an existing Hospital record.
func (r *hospitalRepo) Update(h Hospital) (*Hospital, error) {
	h.UpdatedAt = time.Now()

	// Use named query for easier parameter binding
	query := `
		UPDATE hospitals
		SET name=:name, address=:address, phone_number=:phone_number, email=:email, updated_at=:updated_at
		WHERE hospital_id = :hospital_id
		RETURNING hospital_id, name, address, phone_number, email, created_at, updated_at
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
	return nil, ErrFailedUpdate // Return custom public error
}

// Delete removes a Hospital record by ID.
func (r *hospitalRepo) Delete(id int) error {
	query := `DELETE from hospitals WHERE hospital_id = $1`

	result, err := r.dbCon.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error executing delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Could not get rows affected after delete: %v", err)
		// Continue, as delete was likely successful but we can't confirm.
		return nil
	}

	if rowsAffected == 0 {
		return ErrNotFound // Return custom public error if no row was deleted
	}

	return nil
}

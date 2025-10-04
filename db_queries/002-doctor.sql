CREATE TABLE doctors (
    doctor_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    specialty VARCHAR(100),
    years_experience INT CHECK (years_experience >= 0),
    phone_number VARCHAR(50) UNIQUE,
    email VARCHAR(100) UNIQUE,
    image_url VARCHAR(355),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
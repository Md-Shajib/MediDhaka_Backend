CREATE TABLE hospital_doctor (
    hospital_id INT NOT NULL,
    doctor_id INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    role VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (hospital_id, doctor_id),
    FOREIGN KEY (hospital_id) REFERENCES hospitals(hospital_id) ON DELETE CASCADE,
    FOREIGN KEY (doctor_id) REFERENCES doctors(doctor_id) ON DELETE CASCADE
);
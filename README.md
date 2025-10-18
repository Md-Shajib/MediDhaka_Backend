# Medidhaka

Medidhaka is a RESTful API service for managing hospitals, doctors, and their relationships. It provides endpoints to create, read, update, and delete hospital and doctor records, as well as to manage associations between them.

---

## Table of Contents

- [Features](#features)  
- [Tech Stack](#tech-stack)  
- [Getting Started](#getting-started)  
- [API Endpoints](#api-endpoints)  
- [Project Structure](#project-structure)  
- [Middleware](#middleware)  
- [Database Schema](#database-schema)

---

## Features

- Manage Hospitals: create, retrieve, update, delete hospital records
- Manage Doctors: create, retrieve, update, delete doctor records
- Manage Hospital-Doctor relationships, including assigning roles
- Search functionality for hospitals and doctors by name
- Pagination support for listing endpoints
- CORS middleware to enable cross-origin requests

---

## Tech Stack

- **Go (Golang)** for backend development
- **PostgreSQL** as the primary relational database
- **github.com/jmoiron/sqlx** for SQL database interaction
- **Gorilla Mux** for HTTP routing
- Middleware management for global and route-level middleware support
- JSON-based REST API responses

---

## Getting Started

### Prerequisites

- Go 1.20+ installed ([download here](https://golang.org/dl/))
- PostgreSQL installed and running
- `medidhaka` database created with required tables (see [Database Schema](#database-schema) below)

### Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/Md-Shajib/MediDhaka_Backend.git
   cd medidhaka
2. Configure your PostgreSQL connection string in db/connection.go:
    ```bash
    func GetConnectionString() string {
        return "user=db_userName password=db_password host=your_host port=db_port dbname=your_db_name sslmode=disable"
    }
3. Run migrations or manually create tables
4. Build and run the server:
   ```bash
   go build -o medidhaka ./cmd/server
    ./medidhaka
5. API server listens on port 8080 by default. Otherwise you have to define the port. Here I am using .env file to set the port.
---

## API Endpoints

### i. Hospitals

| Method | Endpoint          | Description                             |
| ------ | ----------------- | --------------------------------------- |
| POST   | `/hospitals`      | Create a new hospital                   |
| GET    | `/hospitals`      | List hospitals with search & pagination |
| GET    | `/hospitals/{id}` | Get hospital by ID                      |
| PUT    | `/hospitals/{id}` | Update hospital by ID                   |
| DELETE | `/hospitals/{id}` | Delete hospital by ID                   |

### ii. Doctors

| Method | Endpoint        | Description                           |
| ------ | --------------- | ------------------------------------- |
| POST   | `/doctors`      | Create a new doctor                   |
| GET    | `/doctors`      | List doctors with search & pagination |
| GET    | `/doctors/{id}` | Get doctor by ID                      |
| PUT    | `/doctors/{id}` | Update doctor by ID                   |
| DELETE | `/doctors/{id}` | Delete doctor by ID                   |

### iii. Hospital-Doctor Relationship

| Method | Endpoint                                     | Description                         |
| ------ | -------------------------------------------- | ----------------------------------- |
| POST   | `/hospital-doctor`                           | Assign a doctor to a hospital       |
| GET    | `/hospital-doctor/{hospital_id}`             | List doctors assigned to a hospital |
| DELETE | `/hospital-doctor/{hospital_id}/{doctor_id}` | Remove doctor-hospital association  |

### iv. Global Search

| Method | Endpoint  | Description                          |
| ------ | --------- | ------------------------------------ |
| GET    | `/search` | Search doctors and hospitals by name |

---


## Project Structure
``` bash
medidhaka/
├── cmd/
│   └── serve.go               # Entry point of the application
├── config
│   └── config.go              # Database Configuration
├── db/
│   └── connection.go          # Database connection setup
├── repo/                      # Repository layer for DB operations
│   ├── doctor_repo.go
│   ├── hospital_repo.go
│   ├── hospital_doctor_repo.go
│   └── ...
├── rest/                      # REST API handlers and routing
│   ├── handlers/
│   │   ├── doctors.go              
│   │   ├── hospitals.go
│   │   ├── hospital_doctor.go          
│   │   └── search_handler.go
│   ├── middlewares/
│   │   ├── cors.go           # CORS middleware 
│   │   ├── ...               
│   │   └── manager.go        # Middleware manager for chaining
│   ├── routes.go
│   └── server.go             # Repository & Route Initialize
└── util/
    └── send_data.go          # Utility functions for response formatting
```
---

## Middleware

- CORS Middleware: Allows cross-origin resource sharing by setting appropriate headers.

- Middleware Manager: Supports registering global and route-specific middlewares with clean chaining.

## Database Schema
``` bash
CREATE TABLE hospitals (
  hospital_id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  address TEXT NOT NULL,
  phone_number VARCHAR(50),
  email VARCHAR(255),
  image_url TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE doctors (
  doctor_id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  specialty VARCHAR(255),
  years_experience INT,
  phone_number VARCHAR(50),
  email VARCHAR(255),
  image_url TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE hospital_doctor (
  hospital_id INT REFERENCES hospitals(hospital_id) ON DELETE CASCADE,
  doctor_id INT REFERENCES doctors(doctor_id) ON DELETE CASCADE,
  role VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (hospital_id, doctor_id)
);
```
---

---
---
<h1 align="center"> Thank You 🌹</h1>
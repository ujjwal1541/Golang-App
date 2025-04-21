# Healthcare Management API

A Golang web application with receptionist and doctor portals for patient management.

## Features

### Authentication & Users
- Single login for both receptionists and doctors with JWT
- Role-based access control
- Secure password hashing

### Receptionist Portal
- Register new patients
- View, update, and delete patient records
- Search for patients

### Doctor Portal
- View patient details
- Update patient medical information

## Technology Stack

- **Backend**: Golang with Gin web framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT
- **Documentation**: Swagger
- **Testing**: Unit tests with testify

## Project Structure

```
.
├── cmd
│   └── api
│       └── main.go     # Application entry point
├── config
│   └── config.go       # Configuration handling
├── internal
│   ├── handlers        # HTTP handlers
│   ├── middleware      # Middleware functions
│   ├── models          # Data models
│   ├── repositories    # Database access layer
│   └── services        # Business logic layer
├── docs                # Documentation
├── go.mod              # Go modules file
├── go.sum              # Go modules checksums
└── README.md           # This file
```

## API Endpoints

### Authentication
- `POST /api/v1/login` - Login for both doctor and receptionist

### Users
- `POST /api/v1/users` - Create a new user (requires authentication)

### Patients (Receptionist Access)
- `POST /api/v1/patients` - Register a new patient
- `GET /api/v1/patients` - Get all patients with pagination
- `GET /api/v1/patients/:id` - Get a specific patient
- `PUT /api/v1/patients/:id` - Update patient information
- `DELETE /api/v1/patients/:id` - Delete a patient

### Patients (Doctor Access)
- `GET /api/v1/doctor/patients` - Get all patients with pagination
- `GET /api/v1/doctor/patients/:id` - Get a specific patient
- `PUT /api/v1/doctor/patients/:id/medical` - Update patient medical information

## Setup and Installation

### Prerequisites
- Go 1.21 or higher
- PostgreSQL
- Git

### Installation

1. Clone the repository
   ```
   git clone <repository-url>
   cd healthcare-app
   ```

2. Set up environment variables
   ```
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=healthcare
   export JWT_SECRET=your-256-bit-secret
   export SERVER_PORT=8080
   ```

3. Run the application
   ```
   go run cmd/api/main.go
   ```

4. Access the API documentation
   ```
   http://localhost:8080/swagger/index.html
   ```

## Testing

Run the tests with:
```
go test ./internal/...
```

## Design Patterns Used

- **Repository Pattern**: Separation of database access logic
- **Service Layer**: Business logic separate from HTTP handlers
- **Dependency Injection**: For better testability and modularity
- **Clean Architecture**: Separation of concerns with clear layers

## Database Schema

- **Users**: Store user credentials and roles
- **Patients**: Store patient information with medical details

## Future Improvements

- Add more features like appointments scheduling
- Implement a frontend client
- Add more comprehensive logs
- Implement more validation rules
- Add database transaction support 
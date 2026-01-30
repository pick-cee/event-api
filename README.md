# Event Management API

A REST API for managing events and registrations built with Go, Gin, PostgreSQL, and GORM.

## Features

- ✅ User authentication (JWT)
- ✅ Create, read, update, delete events
- ✅ Event registration system
- ✅ Authorization (users can only modify their own events)
- ✅ View event attendees

## Tech Stack

- **Framework:** Gin
- **Database:** PostgreSQL
- **ORM:** GORM
- **Authentication:** JWT
- **Password Hashing:** bcrypt

## Project Structure

```
event-api/
├── cmd/api/              # Application entry point
├── internal/
│   ├── models/           # Database models
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Middleware (auth, etc.)
│   ├── database/         # Database connection & migrations
│   └── config/           # Configuration
├── .env                  # Environment variables
└── README.md
```

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 14+

### Installation

1. Clone the repository

```bash
git clone https://github.com/pick-cee/event-api.git
cd event-api
```

2. Install dependencies

```bash
go mod download
```

3. Create `.env` file

```bash
cp .env.example .env
```

4. Update `.env` with your database credentials

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=event_api
JWT_SECRET=your-secret-key
```

5. Create database

```bash
createdb event_api
```

6. Run the application

```bash
go run cmd/api/*.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

| Method | Endpoint              | Description       | Auth Required |
| ------ | --------------------- | ----------------- | ------------- |
| POST   | `/api/v1/auth/signup` | Register new user | No            |
| POST   | `/api/v1/auth/login`  | Login user        | No            |

### Events

| Method | Endpoint             | Description                 | Auth Required |
| ------ | -------------------- | --------------------------- | ------------- |
| GET    | `/api/v1/events`     | List all events             | No            |
| GET    | `/api/v1/events/:id` | Get single event            | No            |
| POST   | `/api/v1/events`     | Create event                | Yes           |
| PUT    | `/api/v1/events/:id` | Update event (creator only) | Yes           |
| DELETE | `/api/v1/events/:id` | Delete event (creator only) | Yes           |

### Registrations

| Method | Endpoint                       | Description          | Auth Required |
| ------ | ------------------------------ | -------------------- | ------------- |
| POST   | `/api/v1/events/:id/register`  | Register for event   | Yes           |
| DELETE | `/api/v1/events/:id/register`  | Cancel registration  | Yes           |
| GET    | `/api/v1/events/:id/attendees` | Get event attendees  | No            |
| GET    | `/api/v1/my-registrations`     | Get my registrations | Yes           |

## Usage Examples

### 1. Signup

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Response:**

```json
{
	"user": {
		"id": 1,
		"name": "John Doe",
		"email": "john@example.com"
	},
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Create Event

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Go Workshop",
    "description": "Learn Go programming",
    "location": "Lagos, Nigeria",
    "date_time": "2025-02-15T10:00:00Z"
  }'
```

### 4. List Events

```bash
curl http://localhost:8080/api/v1/events
```

### 5. Register for Event

```bash
curl -X POST http://localhost:8080/api/v1/events/1/register \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 6. Get My Registrations

```bash
curl http://localhost:8080/api/v1/my-registrations \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 7. Update Event (Creator Only)

```bash
curl -X PUT http://localhost:8080/api/v1/events/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Advanced Go Workshop",
    "location": "Abuja, Nigeria"
  }'
```

### 8. Delete Event (Creator Only)

```bash
curl -X DELETE http://localhost:8080/api/v1/events/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Database Schema

### Users

- `id` (Primary Key)
- `name`
- `email` (Unique)
- `password` (Hashed)
- `created_at`
- `updated_at`

### Events

- `id` (Primary Key)
- `title`
- `description`
- `location`
- `date_time`
- `creator_id` (Foreign Key → Users)
- `created_at`
- `updated_at`

### Registrations

- `id` (Primary Key)
- `user_id` (Foreign Key → Users)
- `event_id` (Foreign Key → Events)
- `created_at`

## License

MIT

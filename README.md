````markdown
# Event Management API

A production-ready REST API for managing events and registrations built with Go, Gin, PostgreSQL, GORM, Redis, and automated email notifications.

## Features

- ✅ User authentication & authorization (JWT)
- ✅ Create, read, update, delete events
- ✅ Event registration system
- ✅ Authorization (users can only modify their own events)
- ✅ View event attendees
- ✅ Email notifications (Novu integration)
  - Welcome emails on signup
  - Registration confirmation emails
  - Event reminder emails (24 hours & 1 hour before)
- ✅ Redis caching for performance
- ✅ Automated cron jobs for event reminders
- ✅ Pagination support
- ✅ Rate limiting ready
- ✅ CORS support
- ✅ Graceful shutdown

## Tech Stack

- **Framework:** Gin
- **Database:** PostgreSQL
- **ORM:** GORM
- **Cache:** Redis
- **Authentication:** JWT
- **Password Hashing:** bcrypt
- **Email Service:** Novu
- **Job Scheduler:** gocron

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Novu account (for email notifications)

### Installation

1. Clone the repository

```bash
git clone https://github.com/pick-cee/events-api.git
cd events-api
```
````

2. Install dependencies

```bash
go mod download
```

3. Create `.env` file

```bash
cp .env.example .env
```

4. Update `.env` with your credentials

```env
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=event_api

# JWT
JWT_SECRET=your-super-secret-key-change-in-production

# Redis
REDIS_URL=redis://localhost:6379

# Novu (Email Service)
NOVU_SECRET_KEY=your-novu-secret-key
```

5. Start PostgreSQL and Redis

```bash
# PostgreSQL
createdb event_api

# Redis (if not already running)
redis-server
```

6. Run the application

```bash
go run cmd/api/*.go
```

The server will start on `http://localhost:8080`

**You should see:**

```
✅ Configuration loaded
✅ Connected to database
✅ Migrations completed
✅ Connected to Redis
✅ Scheduler started
  - 24h reminders: Every 1 hour
  - 1h reminders: Every 10 minutes
🚀 Server running on port 8080
```

## API Endpoints

### Authentication

| Method | Endpoint              | Description       | Auth Required |
| ------ | --------------------- | ----------------- | ------------- |
| POST   | `/api/v1/auth/signup` | Register new user | No            |
| POST   | `/api/v1/auth/login`  | Login user        | No            |

### Events

| Method | Endpoint             | Description                 | Auth Required |
| ------ | -------------------- | --------------------------- | ------------- |
| GET    | `/api/v1/events`     | List all events (paginated) | No            |
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

## Cron Jobs

The API runs automated jobs for event reminders:

| Job               | Schedule         | Description                                 |
| ----------------- | ---------------- | ------------------------------------------- |
| 24-hour reminders | Every 1 hour     | Sends reminders for events happening in 24h |
| 1-hour reminders  | Every 10 minutes | Sends reminders for events happening in 1h  |

Jobs use Redis to prevent duplicate emails.

## Database Schema

### Users

- `id` (Primary Key)
- `name`
- `email` (Unique)
- `password` (Hashed with bcrypt)
- `created_at`
- `updated_at`
- `deleted_at` (Soft delete)

### Events

- `id` (Primary Key)
- `title`
- `description`
- `location`
- `date_time`
- `creator_id` (Foreign Key → Users)
- `created_at`
- `updated_at`
- `deleted_at` (Soft delete)

### Registrations

- `id` (Primary Key)
- `user_id` (Foreign Key → Users)
- `event_id` (Foreign Key → Events)
- `created_at`
- `deleted_at` (Soft delete)

## Caching

Redis is used for:

- Preventing duplicate reminder emails
- Session management (future feature)
- Rate limiting (future feature)

## License

MIT

## Acknowledgments

- [Gin](https://gin-gonic.com/) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [Novu](https://novu.co/) - Notification infrastructure
- [gocron](https://github.com/go-co-op/gocron) - Job scheduler
- [go-redis](https://github.com/redis/go-redis) - Redis client

````

# GoSplit

A REST API for splitting expenses between groups — a Splitwise clone built as a learning project.

## Stack

- **Language:** Go (standard library only — `net/http`, `database/sql`)
- **Database:** PostgreSQL
- **Auth:** JWT (`golang-jwt/jwt`)
- **Other:** `lib/pq` driver, `google/uuid`, `godotenv`
- **Infra:** Docker, Docker Compose

No Gin. No GORM.

## Features

- User registration and login with JWT authentication
- Create and manage groups
- Add and remove group members
- Record expenses with per-user splits
- Record settlements between users
- Calculate net balances per user in a group

## Project Structure

```
cmd/
  main.go                  # Entry point
internal/
  auth/
    handler.go             # POST /api/auth/register, POST /api/auth/login
    service.go
  groups/
    handler.go             # CRUD + member management
    service.go
  expenses/
    handler.go             # CRUD + splits
    service.go
  settlements/
    handler.go             # Create + list settlements
    service.go
  balances/
    handler.go             # GET /api/groups/{id}/balances
    service.go
migrations/
  001_init.sql             # All 6 tables
pkg/
  database/
    postgres.go            # DB connection
  middleware/
    auth.go                # JWT middleware + GetUserID helper
  models/
    models.go              # Shared structs
```

## Getting Started

### Prerequisites

- Docker + Docker Compose

### Run

```bash
docker compose up --build
```

The API will be available at `http://localhost:8080`.

### Environment Variables

Create a `.env` file in the project root:

```env
POSTGRES_USER=gosplit
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=gosplit_db
POSTGRES_PORT=5432
APP_PORT=8080
JWT_SECRET=your_jwt_secret
DATABASE_URL=postgres://gosplit:yourpassword@localhost:5432/gosplit_db?sslmode=disable
```

## API Reference

### Auth

| Method | Route | Description | Auth |
|--------|-------|-------------|------|
| POST | `/api/auth/register` | Register a new user | ❌ |
| POST | `/api/auth/login` | Login and get JWT token | ❌ |

### Groups

| Method | Route | Description | Auth |
|--------|-------|-------------|------|
| POST | `/api/groups` | Create a group | ✅ |
| GET | `/api/groups` | List user's groups | ✅ |
| GET | `/api/groups/{id}` | Get a group | ✅ |
| PUT | `/api/groups/{id}` | Update a group | ✅ |
| DELETE | `/api/groups/{id}` | Delete a group | ✅ |
| POST | `/api/groups/{id}/members` | Add a member | ✅ |
| DELETE | `/api/groups/{id}/members/{user_id}` | Remove a member | ✅ |

### Expenses

| Method | Route | Description | Auth |
|--------|-------|-------------|------|
| POST | `/api/groups/{id}/expenses` | Create an expense | ✅ |
| GET | `/api/groups/{id}/expenses` | List expenses in a group | ✅ |
| GET | `/api/groups/{id}/expenses/{expenseId}` | Get an expense | ✅ |
| DELETE | `/api/groups/{id}/expenses/{expenseId}` | Delete an expense | ✅ |

### Settlements

| Method | Route | Description | Auth |
|--------|-------|-------------|------|
| POST | `/api/groups/{id}/settlements` | Record a settlement | ✅ |
| GET | `/api/groups/{id}/settlements` | List settlements in a group | ✅ |

### Balances

| Method | Route | Description | Auth |
|--------|-------|-------------|------|
| GET | `/api/groups/{id}/balances` | Get net balances for all users in a group | ✅ |

## Balance Calculation

A user's balance in a group is calculated as:

```
balance = expenses paid by user
        - splits assigned to user
        + settlements received
        - settlements paid
```

A **positive** balance means the user is owed money.
A **negative** balance means the user owes money.

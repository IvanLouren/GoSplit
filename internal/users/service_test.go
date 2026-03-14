package users_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/IvanLouren/GoSplit/internal/users"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("gosplit_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	defer pgContainer.Terminate(ctx)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %s", err)
	}
	defer testDB.Close()

	if err := runMigrations(testDB); err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}

	os.Exit(m.Run())
}

func runMigrations(db *sql.DB) error {
	migration, err := os.ReadFile("../../migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration: %w", err)
	}
	_, err = db.Exec(string(migration))
	return err
}

func TestGetMe(t *testing.T) {
	var userID string
	err := testDB.QueryRow(`INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`,
		"User 1", "user1@test.com", "hashedpassword").Scan(&userID)
	if err != nil {
		t.Fatalf("failed to insert user: %s", err)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		t.Fatalf("failed to parse userID: %s", err)
	}

	service := users.NewService(testDB)
	user, err := service.GetMe(parsedUserID)
	if err != nil {
		t.Fatalf("failed to get user: %s", err)
	}

	if user.ID != parsedUserID {
		t.Errorf("expected user ID %s, got %s", parsedUserID, user.ID)
	}

	if user.Name != "User 1" {
		t.Errorf("expected user name 'User 1', got %s", user.Name)
	}

	if user.Email != "user1@test.com" {
		t.Errorf("expected user email 'user1@test.com', got %s", user.Email)
	}
}

func TestUpdateMe(t *testing.T) {
	var userID string
	err := testDB.QueryRow(
		`INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`,
		"User 2", "user2@test.com", "hashedpassword",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to insert user: %s", err)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		t.Fatalf("failed to parse userID: %s", err)
	}

	service := users.NewService(testDB)
	updated, err := service.UpdateMe(parsedUserID, "User 2 Updated")
	if err != nil {
		t.Fatalf("failed to update me: %s", err)
	}

	if updated.ID != parsedUserID {
		t.Errorf("expected id %s, got %s", parsedUserID, updated.ID)
	}
	if updated.Name != "User 2 Updated" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}
	if updated.Email != "user2@test.com" {
		t.Errorf("expected email unchanged, got %s", updated.Email)
	}
}

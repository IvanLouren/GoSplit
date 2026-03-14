package users

import (
	"database/sql"

	"github.com/IvanLouren/GoSplit/pkg/models"
	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetMe(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		`SELECT id, name, email, password, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) UpdateMe(userID uuid.UUID, name string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		`UPDATE users SET name = $1 WHERE id = $2 RETURNING id, name, email, password, created_at`,
		name, userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

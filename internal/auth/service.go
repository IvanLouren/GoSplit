package auth

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/IvanLouren/GoSplit/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Register(name, email, password string) (*models.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New()
	createdAt := time.Now()
	_, err = s.db.Exec(`INSERT INTO users (id, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5)`,
		id, name, email, string(hash), createdAt)

	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  string(hash),
		CreatedAt: createdAt,
	}, nil
}

func (s *Service) Login(email, password string) (string, error) {

	var user models.User
	err := s.db.QueryRow(`SELECT id, name, email, password, created_at FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return "", errors.New("invalid credentials")
	}
	if err != nil {
		return "", err

	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))

}

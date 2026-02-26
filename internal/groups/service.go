package groups

import (
	"database/sql"
	"time"

	"github.com/IvanLouren/GoSplit/pkg/models"
	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateGroup(name string, createdBy uuid.UUID) (*models.Group, error) {
	groupID := uuid.New()
	createdAt := time.Now()

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // rolls back if we don't commit

	_, err = tx.Exec(`INSERT INTO groups (id, name, created_by, created_at) VALUES ($1, $2, $3, $4)`,
		groupID, name, createdBy, createdAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`INSERT INTO group_members (id, group_id, user_id, joined_at) VALUES ($1, $2, $3, $4)`,
		uuid.New(), groupID, createdBy, createdAt)
	if err != nil {
		return nil, err
	}

	// Commit — both inserts succeed or neither does
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Group{ID: groupID, Name: name, CreatedBy: createdBy, CreatedAt: createdAt}, nil
}

func (s *Service) GetGroups(userID uuid.UUID) ([]models.Group, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(`
		SELECT g.id, g.name, g.created_by, g.created_at
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (s *Service) GetGroup(groupID uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := s.db.QueryRow(`SELECT id, name, created_by, created_at FROM groups WHERE id = $1`, groupID).
		Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *Service) UpdateGroup(groupID uuid.UUID, name string) (*models.Group, error) {
	_, err := s.db.Exec(`UPDATE groups SET name = $1 WHERE id = $2`, name, groupID)
	if err != nil {
		return nil, err
	}
	return s.GetGroup(groupID) // now reads after the update is committed
}

func (s *Service) DeleteGroup(groupID uuid.UUID) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`DELETE FROM groups WHERE id = $1`, groupID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) AddMember(groupID, userID uuid.UUID) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO group_members (id, group_id, user_id, joined_at) VALUES ($1, $2, $3, $4)`,
		uuid.New(), groupID, userID, time.Now())
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) RemoveMember(groupID, userID uuid.UUID) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`, groupID, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

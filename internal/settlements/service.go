package settlements

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

func (s *Service) CreateSettlement(groupID, paidBy, paidTo uuid.UUID, amount float64) (models.Settlement, error) {
	var settlement models.Settlement
	err := s.db.QueryRow(`INSERT INTO settlements (group_id, paid_by, paid_to, amount) VALUES
					 ($1, $2, $3, $4) RETURNING id, group_id, paid_by, paid_to, amount, created_at`, groupID, paidBy, paidTo, amount).
		Scan(&settlement.ID, &settlement.GroupID, &settlement.PaidBy, &settlement.PaidTo, &settlement.Amount, &settlement.CreatedAt)
	if err != nil {
		return models.Settlement{}, err
	}

	return settlement, nil
}

func (s *Service) GetSettlements(groupID uuid.UUID) ([]models.Settlement, error) {
	settlements, err := s.db.Query(`SELECT id, group_id, paid_by, paid_to, amount, created_at FROM settlements WHERE group_id = $1`, groupID)
	if err != nil {
		return nil, err
	}
	defer settlements.Close()

	var result []models.Settlement
	for settlements.Next() {
		var settlement models.Settlement
		err = settlements.Scan(&settlement.ID, &settlement.GroupID, &settlement.PaidBy, &settlement.PaidTo, &settlement.Amount, &settlement.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, settlement)
	}

	return result, nil
}

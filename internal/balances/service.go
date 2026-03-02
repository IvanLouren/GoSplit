package balances

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

func (s *Service) GetBalances(groupID uuid.UUID) ([]models.Balance, error) {
	// single SQL query using UNION ALL
	balances, err := s.db.Query(`SELECT user_id, SUM(amount) as balance 
		FROM ( SELECT paid_by AS user_id, amount
		   	   FROM expenses
		       WHERE group_id = $1
				UNION ALL
			   SELECT expense_splits.user_id, -expense_splits.amount
			   FROM expense_splits
			   JOIN expenses ON expenses.id = expense_splits.expense_id
			   WHERE expenses.group_id = $1
			    UNION ALL
			   SELECT paid_by as user_id, -amount
			   FROM settlements
			   WHERE group_id =$1
				UNION ALL
			   SELECT paid_to AS user_id, amount
			   FROM settlements
			   WHERE group_id = $1
		) as entries
		GROUP BY user_id
	 `, groupID)
	if err != nil {
		return nil, err
	}
	defer balances.Close()

	var result []models.Balance
	for balances.Next() {
		var balance models.Balance
		err = balances.Scan(&balance.UserID, &balance.Balance)
		if err != nil {
			return nil, err
		}
		result = append(result, balance)
	}

	return result, nil
}

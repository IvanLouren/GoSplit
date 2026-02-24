package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Group struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type GroupMember struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	GroupID  uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

type Expense struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	GroupID     uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	PaidBy      uuid.UUID `gorm:"type:uuid;not null" json:"paid_by"`
	Description string    `gorm:"not null" json:"description"`
	Amount      float64   `gorm:"not null" json:"amount"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type ExpenseSplit struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ExpenseID uuid.UUID `gorm:"type:uuid;not null" json:"expense_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Amount    float64   `gorm:"not null" json:"amount"`
}

type Settlement struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	PaidBy    uuid.UUID `gorm:"type:uuid;not null" json:"paid_by"`
	PaidTo    uuid.UUID `gorm:"type:uuid;not null" json:"paid_to"`
	Amount    float64   `gorm:"not null" json:"amount"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

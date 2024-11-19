package entities

import (
	"time"
)

type Transaction struct {
	ID        int        `db:"id" json:"id"`
	UserID    int        `db:"user_id" json:"userId"`
	BankID    int        `db:"bank_id" json:"bankId"`
	Amount    float64    `db:"amount" json:"amount"`
	Type      string     `db:"type" json:"type"`
	Reference string     `db:"reference" json:"reference"`
	Status    string     `db:"status" json:"status"`
	ExpiredAt time.Time  `db:"expired_at" json:"expiredAt"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

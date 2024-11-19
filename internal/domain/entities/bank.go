package entities

import (
	"time"
)

type Bank struct {
	ID            int        `db:"id" json:"id"`
	UserID        int        `db:"user_id" json:"userId"`
	AccountName   string     `db:"account_name" json:"accountName"`
	AccountNumber string     `db:"account_number" json:"accountNumber"`
	BankName      string     `db:"bank_name" json:"bankName"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deletedAt"`
}

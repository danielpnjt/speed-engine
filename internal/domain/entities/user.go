package entities

import (
	"time"
)

type User struct {
	ID        int        `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Password  string     `db:"password" json:"-"`
	Email     string     `db:"email" json:"email"`
	Name      string     `db:"name" json:"name"`
	Balance   float64    `db:"balance" json:"balance"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

package models

import "time"

type User struct {
	ID        int       `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	Avatar    string    `db:"avatar,omitempty"`
	CreatedAt time.Time `db:"created_at"`
}

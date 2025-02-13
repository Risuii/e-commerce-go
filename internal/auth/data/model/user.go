package model

type User struct {
	Uuid          string `db:"uuid"`
	Email         string `db:"email"`
	EmailVerified bool   `db:"email_verified"`
	Password      string `db:"password"`
	Username      string `db:"username"`
	LastLogin     string `db:"last_login"`
	CreatedAt     string `db:"created_at"`
	UpdatedAt     string `db:"updated_at"`
	DeletedAt     string `db:"deleted_at"`
}

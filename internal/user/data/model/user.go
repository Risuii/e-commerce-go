package model

type User struct {
	Uuid          string `db:"uuid" gorm:"column:uuid"`
	Email         string `db:"email" gorm:"column:email"`
	EmailVerified bool   `db:"email_verified" gorm:"column:email_verified"`
	Password      string `db:"password" gorm:"column:password"`
	Username      string `db:"username" gorm:"column:username"`
	LastLogin     string `db:"last_login" gorm:"column:last_login"`
	CreatedAt     string `db:"created_at" gorm:"column:created_at"`
	UpdatedAt     string `db:"updated_at" gorm:"column:updated_at"`
	DeletedAt     string `db:"deleted_at" gorm:"column:deleted_at"`
}

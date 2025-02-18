package model

type Store struct {
	StoreID     string `db:"store_id" gorm:"column:store_id"`
	StoreName   string `db:"store_name" gorm:"column:store_name"`
	Description string `db:"description" gorm:"column:description"`
	UserID      string `db:"user_id" gorm:"column:user_id"`
	CreatedAt   string `db:"created_at" gorm:"column:created_at"`
	UpdatedAt   string `db:"updated_at" gorm:"column:updated_at"`
}

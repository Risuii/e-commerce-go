package dto

type ProductParam struct {
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

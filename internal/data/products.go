package data

import "database/sql"

type Product struct {
	ProductID   int
	ProductName *string
	Price       uint64
	Rating      uint8
	Image       *string
}

type ProductModel struct {
	DB *sql.DB
}

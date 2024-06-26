package data

import (
	"errors"

	"github.com/vaidik-bajpai/ecommerce-api/internal/prisma/db"
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrDuplicateEmail   = errors.New("error duplicate email")
	ErrDuplicatePhoneNo = errors.New("error duplicate phone no")
	ErrMultipleCarts    = errors.New("error user cannot have more than one cart")
)

type Models struct {
	Users    UserModel
	Products ProductModel
	Carts    CartModel
}

func NewModels(db *db.PrismaClient) Models {
	return Models{
		Users:    UserModel{DB: db},
		Products: ProductModel{DB: db},
		Carts:    CartModel{DB: db},
	}
}

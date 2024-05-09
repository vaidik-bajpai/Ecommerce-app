package data

import (
	"errors"

	"github.com/vaidik-bajpai/ecommerce-api/prisma/db"
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrDuplicateEmail   = errors.New("error duplicate email")
	ErrDuplicatePhoneNo = errors.New("error duplicate phone no")
)

type Models struct {
	Users    UserModel
	Products ProductModel
}

func NewModels(db *db.PrismaClient) Models {
	return Models{
		Users:    UserModel{DB: db},
		Products: ProductModel{DB: db},
	}
}

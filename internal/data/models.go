package data

import (
	"database/sql"
	"errors"
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

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Products: ProductModel{DB: db},
	}
}

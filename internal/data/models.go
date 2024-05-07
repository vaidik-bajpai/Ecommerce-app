package data

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicateEmail   = errors.New("error duplicate email")
	ErrDuplicatePhoneNo = errors.New("error duplicate phone no")
)

type Models struct {
	Users UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}

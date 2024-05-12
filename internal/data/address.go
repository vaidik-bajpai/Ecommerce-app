package data

import (
	"context"
	"time"

	"github.com/vaidik-bajpai/ecommerce-api/internal/prisma/db"
)

type Address struct {
	ID      int
	House   *string
	Street  *string
	City    *string
	Pincode *string
	UserID  int
}

type AddressModel struct {
	DB *db.PrismaClient
}

func (m AddressModel) AddAddress(userId int, address *Address) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newAddress, err := m.DB.Address.CreateOne(
		db.Address.House.Set(*address.House),
		db.Address.Street.Set(*address.Street),
		db.Address.City.Set(*address.City),
		db.Address.Pincode.Set(*address.Pincode),
		db.Address.User.Link(
			db.User.ID.Equals(userId),
		),
	).Exec(ctx)

	if err != nil {
		return err
	}

	address.ID = newAddress.ID

	return nil
}

func (m AddressModel) RemoveAddress(addressID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Product.FindUnique(
		db.Product.ID.Equals(addressID),
	).Delete().Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

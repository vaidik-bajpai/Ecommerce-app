package data

import (
	"context"
	"errors"
	"time"

	"github.com/vaidik-bajpai/ecommerce-api/internal/prisma/db"
)

type Cart struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

type CartModel struct {
	DB *db.PrismaClient
}

func (m CartModel) AddToCart(product *Product, userId, productId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Cart.UpsertOne(
		db.Cart.UserID.Equals(userId),
	).Create(
		db.Cart.UserID.Set(userId),
	).Update(
		db.Cart.UserID.Set(userId),
	).Exec(ctx)

	if err != nil {
		return err
	}

	_, err = m.DB.Cart.FindUnique(
		db.Cart.UserID.Equals(userId),
	).Update(
		db.Cart.User.Link(
			db.User.ID.Equals(userId),
		),
	).Exec(ctx)

	/* if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			createErr := m.CreateCart(userId)
			if createErr != nil {
				return createErr
			}
		default:
			return err
		}
	} */

	if err != nil {
		return err
	}

	return nil
}

func (m CartModel) CreateCart(userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Cart.CreateOne(
		db.Cart.User.Link(
			db.User.ID.Equals(userId),
		),
	).Exec(ctx)

	if err != nil {
		infoUnique, isErr := db.IsErrUniqueConstraint(err)

		switch {
		case isErr:
			for _, field := range infoUnique.Fields {
				if field == "userid" {
					return ErrMultipleCarts
				} else {
					return errors.New("unique constraint violated")
				}
			}
		default:
			return err
		}
	}

	return nil
}

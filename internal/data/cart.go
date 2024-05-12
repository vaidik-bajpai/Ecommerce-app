package data

import (
	"context"
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

func (m CartModel) AddToCart(product *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userCart, err := m.DB.Cart.CreateOne(
		db.Cart.User.Link(
			db.User.ID.Equals(cart.UserID),
		),
		db.Cart.UserID.Set(cart.UserID),
	).Exec(ctx)

	if err != nil {
		return err
	}

	cart, err := m.DB.Cart.FindUnique(
		prisma.Cart.ID.Equals(cartID),
	).Update(
		prisma.Cart.Products.Connect(
			&prisma.ProductWhereUniqueInput{
				ID: &productID,
			},
		),
	).Exec(ctx)

	return nil

}

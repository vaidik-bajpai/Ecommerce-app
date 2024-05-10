package data

import (
	"time"

	"github.com/vaidik-bajpai/ecommerce-api/internal/prisma/db"
)

type Product struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Price     uint64    `json:"price"`
	Rating    uint8     `json:"rating"`
	Image     *string   `json:"image"`
}

type ProductModel struct {
	DB *db.PrismaClient
}

/* func (m ProductModel) AddProduct(product *Product) error {
	query := `
		INSERT INTO products (name, price, rating, image)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{product.Name, product.Price, product.Rating, product.Image}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&product.ID,
		&product.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m ProductModel) RemoveProduct(productID int64) error {
	query := `
		DELETE FROM products WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, productID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m ProductModel) Get(productID int64) (*Product, error) {
	query := `
		SELECT id, created_at, name, price, rating, image
		FROM products
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var product Product

	err := m.DB.QueryRowContext(ctx, query, productID).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.Name,
		&product.Price,
		&product.Rating,
		&product.Image,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (m ProductModel) GetAll(productName string, price int, filters Filters) ([]*Product, Metadata, error) {
	query := `
		SELECT count(*) OVER(), id, created_at, name, price, rating, image
		FROM products
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (price <= $2)
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, productName, price)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	products := []*Product{}

	for rows.Next() {
		var product Product

		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.CreatedAt,
			&product.Name,
			&product.Price,
			&product.Rating,
			&product.Image,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return products, metadata, nil
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "product_name", "must be provided")
	v.Check(len(product.Name) >= 3, "product_name", "must contains atleast 3 bytes")
	v.Check(len(product.Name) <= 30, "product_name", "must not contains more than 30 bytes")

	v.Check(product.Price > 0, "product_price", "must be a positive number")

	v.Check(int(product.Rating) > 0, "product_rating", "must be greater than 0")
	v.Check(int(product.Rating) < 5, "product_rating", "must not be greater than 5")

	v.Check(validator.Matches(*product.Image, validator.LinkRX), "image", "must be a valid link")
}
*/

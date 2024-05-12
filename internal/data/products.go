package data

import (
	"context"
	"errors"
	"time"

	"github.com/vaidik-bajpai/ecommerce-api/internal/prisma/db"
	"github.com/vaidik-bajpai/ecommerce-api/internal/validator"
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

func (m ProductModel) AddProduct(product *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newProduct, err := m.DB.Product.CreateOne(
		db.Product.Name.Set(product.Name),
		db.Product.Price.Set(int(product.Price)),
		db.Product.Rating.Set(int(product.Rating)),
		db.Product.Image.Set(*product.Image),
	).Exec(ctx)
	if err != nil {
		return err
	}

	createdAt, ok := newProduct.CreatedAt()
	if !ok {
		return errors.New("error accessing created_at of product")
	}

	product.ID = int64(newProduct.ID)
	product.CreatedAt = createdAt

	return nil
}

func (m ProductModel) RemoveProduct(productID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Product.FindUnique(
		db.Product.ID.Equals(int(productID)),
	).Delete().Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (m ProductModel) Get(productID int64) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newProduct, err := m.DB.Product.FindUnique(
		db.Product.ID.Equals(int(productID)),
	).Exec(ctx)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	createdAt, ok := newProduct.CreatedAt()
	if !ok {
		return nil, errors.New("error accessing created at")
	}

	image, ok := newProduct.Image()
	if !ok {
		return nil, errors.New("error accessing image")
	}

	product := Product{
		ID:        int64(newProduct.ID),
		Name:      newProduct.Name,
		Price:     uint64(newProduct.Price),
		Rating:    uint8(newProduct.Rating),
		CreatedAt: createdAt,
		Image:     &image,
	}

	return &product, nil
}

func (m ProductModel) GetAll(productName string, price int, filters Filters) ([]*Product, Metadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var total db.RawInt
	err := m.DB.Prisma.QueryRaw(`SELECT count(*) as total FROM "Product"`).Exec(ctx, total)
	if err != nil {
		return nil, Metadata{}, err
	}

	var records []db.ProductModel

	switch filters.sortColumn() {
	case "id":
		records, err = m.DB.Product.FindMany().Take(filters.limit()).Skip(filters.offset()).OrderBy(
			db.Product.Price.Lte(price),
			db.Product.ID.Order(filters.sortDirection()),
		).Exec(ctx)
	case "price":
		records, err = m.DB.Product.FindMany().Take(filters.limit()).Skip(filters.offset()).OrderBy(
			db.Product.Price.Order(filters.sortDirection()),
		).Exec(ctx)
	case "rating":
		records, err = m.DB.Product.FindMany().Take(filters.limit()).Skip(filters.offset()).OrderBy(
			db.Product.Price.Lte(price),
			db.Product.Rating.Order(filters.sortDirection()),
		).Exec(ctx)
	case "name":
		records, err = m.DB.Product.FindMany().Take(filters.limit()).Skip(filters.offset()).OrderBy(
			db.Product.Price.Lte(price),
			db.Product.Name.Order(filters.sortDirection()),
		).Exec(ctx)
	}

	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			return nil, Metadata{}, ErrRecordNotFound
		default:
			return nil, Metadata{}, err
		}
	}

	var products []*Product

	for _, product := range records {
		createdAt, ok := product.CreatedAt()
		if !ok {
			return nil, Metadata{}, errors.New("error accessing created at")
		}

		image, ok := product.Image()
		if !ok {
			return nil, Metadata{}, errors.New("error accessing image")
		}
		record := Product{
			ID:        int64(product.ID),
			CreatedAt: createdAt,
			Name:      productName,
			Price:     uint64(product.Price),
			Rating:    uint8(product.Rating),
			Image:     &image,
		}

		products = append(products, &record)
	}

	metadata := calculateMetadata(int(total), filters.Page, filters.PageSize)
	return products, metadata, err
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

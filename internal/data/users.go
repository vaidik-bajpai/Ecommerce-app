package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vaidik-bajpai/ecommerce-api/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	FirstName    *string   `json:"firstname"`
	LastName     *string   `json:"lastname"`
	Password     password  `json:"-"`
	Email        *string   `json:"email"`
	Phone        *string   `json:"phone"`
	Token        *string   `json:"-"`
	RefreshToken *string   `json:"-"`
	Version      int       `json:"version"`
	/* UpdatedAt      time.Time
	UserID         string
	UserCart       []ProductUser
	AddressDetails []Address
	OrderStatus    []Order */
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must contain atleast 8 bytes")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes")
}

func ValidateFirstName(v *validator.Validator, firstName string) {
	v.Check(firstName != "", "firstname", "must be provided")
	v.Check(len(firstName) <= 30, "firstName", "must not be more than 30 bytes")
}

func ValidateLastName(v *validator.Validator, lastName string) {
	v.Check(lastName != "", "lastname", "must be provided")
	v.Check(len(lastName) <= 30, "lastName", "must not be more than 30 bytes")
}

func ValidatePhone(v *validator.Validator, phone string) {
	v.Check(phone != "", "phone", "must be provided")
	v.Check(len(phone) == 10, "phone", "must contain 10 digits")
	for i := range phone {
		if phone[i] > '9' || phone[i] < '0' {
			v.AddErrors("phone", "invalid phone number")
			return
		}
	}
}

func ValidateUser(v *validator.Validator, user *User) {
	ValidateFirstName(v, *user.FirstName)
	ValidateLastName(v, *user.LastName)

	ValidateEmail(v, *user.Email)
	ValidatePhone(v, *user.Phone)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash")
	}
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (firstname, lastname, email, password_hash, phone, token, refresh_token)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{user.FirstName, user.LastName, user.Email, user.Password.hash, user.Phone, user.Token, user.RefreshToken}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_phone_key"`:
			return ErrDuplicatePhoneNo
		default:
			return err
		}
	}

	return nil
}

type Address struct {
	ID      int
	House   *string
	Street  *string
	City    *string
	Pincode *string
	UserID  int
}

type ProductUser struct {
	ProductID   int
	ProductName *string
	Price       uint64
	Rating      uint8
	Image       *string
}

type Order struct {
	ID            int
	OrderCart     []ProductUser
	OrderedAt     time.Time
	Price         uint64
	Discount      *int
	PaymentMethod Payment
	UserID        int
}

type Payment struct {
	Digital bool
	COD     bool
}

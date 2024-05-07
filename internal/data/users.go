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

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FirstName *string   `json:"firstname"`
	LastName  *string   `json:"lastname"`
	Password  password  `json:"-"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone"`
	Version   int       `json:"version"`
	/* UpdatedAt      time.Time
	UserID         string
	UserCart       []ProductUser
	AddressDetails []Address
	OrderStatus    []Order */
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
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
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{user.FirstName, user.LastName, user.Email, user.Password.hash, user.Phone}

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

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, firstname, lastname, email, password_hash, phone
		FROM users
		WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password.hash,
		&user.Phone,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Get(userID int64) (*User, error) {
	query := `
		SELECT id, created_at, firstname, lastname, email, hash_password, phone
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password.hash,
		&user.Phone,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

type ProductUser struct {
	ProductID   int
	ProductName *string
	Price       uint64
	Rating      uint8
	Image       *string
}

type Order struct {
	ID            int64
	OrderedAt     time.Time
	Price         uint64
	Discount      *int
	PaymentMethod Payment
	UserID        int64
	AddressID     int64
}

type Payment struct {
	Digital bool
	COD     bool
}

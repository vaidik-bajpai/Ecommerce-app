package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vaidik-bajpai/ecommerce-api/internal/validator"

	"github.com/vaidik-bajpai/ecommerce-api/internal/prisma/db"

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
	Addresses []Address `json:"addresses,omitempty"`
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
	DB *db.PrismaClient
}

func (m UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newUser, err := m.DB.User.CreateOne(
		db.User.FirstName.Set(*user.FirstName),
		db.User.LastName.Set(*user.LastName),
		db.User.Email.Set(*user.Email),
		db.User.Phone.Set(*user.Phone),
		db.User.Password.Set(string(user.Password.hash)),
	).Exec(ctx)

	user.ID = int64(newUser.ID)
	user.Version = newUser.Version
	createdAt, ok := newUser.CreatedAt()
	if !ok {
		return err
	}
	user.CreatedAt = createdAt

	if err != nil {
		return err
	}

	fmt.Println(newUser)

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newUser, err := m.DB.User.FindUnique(
		db.User.Email.Equals(email),
	).With(
		db.User.Addresses.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	fmt.Println(newUser)

	return &user, nil
}

func (m UserModel) Get(userID int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newUser, err := m.DB.User.FindUnique(
		db.User.ID.Equals(int(userID)),
	).With(
		db.User.Addresses.Fetch(),
	).Exec(ctx)

	var user = &User{
		ID:        int64(newUser.ID),
		FirstName: &newUser.FirstName,
		LastName:  &newUser.LastName,
		Email:     &newUser.Email,
		Phone:     &newUser.Phone,
		Version:   newUser.Version,
	}

	for _, dbAddress := range newUser.Addresses() {
		address := Address{
			ID:      dbAddress.ID,
			House:   &dbAddress.House,
			Street:  &dbAddress.Street,
			City:    &dbAddress.City,
			Pincode: &dbAddress.Pincode,
			UserID:  dbAddress.UserID,
		}
		user.Addresses = append(user.Addresses, address)
	}
	user.Password.Set(newUser.Password)

	if err != nil {
		return nil, err
	}

	return user, nil
}

/*
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
*/

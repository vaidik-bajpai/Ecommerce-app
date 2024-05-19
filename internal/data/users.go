package data

import (
	"context"
	"errors"
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
	Cart      []Product `json:"cart"`
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
	v.Check(validator.Matches(phone, validator.PhoneRX), "phone", "must be a valid email address")
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
		db.User.Password.Set(user.Password.hash),
	).Exec(ctx)

	if err != nil {
		infoUnique, isErr := db.IsErrUniqueConstraint(err)

		switch {
		case isErr:
			for _, field := range infoUnique.Fields {
				if field == "email" {
					return ErrDuplicateEmail
				} else if field == "phone" {
					return ErrDuplicatePhoneNo
				} else {
					return errors.New("unique constraint violated")
				}
			}
		default:
			return err
		}
	}

	user.ID = int64(newUser.ID)
	user.Version = newUser.Version
	createdAt, ok := newUser.CreatedAt()
	if !ok {
		return err
	}
	user.CreatedAt = createdAt

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newUser, err := m.DB.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	createdAt, ok := newUser.CreatedAt()
	if !ok {
		return nil, errors.New("something went wrong with our server")
	}

	user = User{
		ID:        int64(newUser.ID),
		FirstName: &newUser.FirstName,
		LastName:  &newUser.LastName,
		Email:     &newUser.Email,
		Phone:     &newUser.Phone,
		Version:   newUser.Version,
		CreatedAt: createdAt,
	}
	user.Password.hash = newUser.Password

	return &user, nil
}

func (m UserModel) Get(userID int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	newUser, err := m.DB.User.FindUnique(
		db.User.ID.Equals(int(userID)),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var user = &User{
		ID:        int64(newUser.ID),
		FirstName: &newUser.FirstName,
		LastName:  &newUser.LastName,
		Email:     &newUser.Email,
		Phone:     &newUser.Phone,
		Version:   newUser.Version,
	}

	user.Password.hash = newUser.Password

	return user, nil
}

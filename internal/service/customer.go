package service

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal"
	"gitlab.com/lokalpay-dev/digital-goods/internal/dto"
	"gitlab.com/lokalpay-dev/digital-goods/internal/entity"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
)

type Customer struct {
	customerRepo internal.CustomerRepositoryItf
	cfg          config.App
}

func NewCustomer(
	customerRepo internal.CustomerRepositoryItf,
	cfg config.App,
) *Customer {
	return &Customer{
		customerRepo: customerRepo,
		cfg:          cfg,
	}
}

func (c *Customer) Registration(payload dto.AuthPayload) (dto.AuthResponse, error) {
	customer, err := c.customerRepo.GetCustomerByEmail(payload.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if customer.Email == payload.Email {
		return dto.AuthResponse{}, errors.New("email has been registered")
	}

	passwordHash := hashAndSalt([]byte(payload.Password))

	customerID, err := c.customerRepo.CreateCustomer(payload.Email, passwordHash)
	if err != nil {
		return dto.AuthResponse{}, nil
	}

	cst := entity.Customer{
		Email: payload.Email,
		ID:    customerID,
		Name:  "",
	}

	token, _ := c.generateJWTToken(cst)

	return dto.AuthResponse{
		Email: payload.Email,
		Token: token,
	}, nil
}

func (c *Customer) Authentication(payload dto.AuthPayload) (dto.AuthResponse, error) {
	customer, err := c.customerRepo.GetCustomerByEmail(payload.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if customer.ID < 1 {
		return dto.AuthResponse{}, errors.New("customer data not found")
	}

	if !comparePasswords(customer.Password, []byte(payload.Password)) {
		return dto.AuthResponse{}, errors.New("invalid password")
	}

	token, err := c.generateJWTToken(customer)
	if err != nil {
		slog.Errorw("failed to generate token", "message", err.Error())
	}

	return dto.AuthResponse{
		Email: customer.Email,
		Token: token,
	}, nil
}

func (c *Customer) GetProfile(id int) (entity.Customer, error) {
	customer, err := c.customerRepo.GetCustomerByID(id)
	if err != nil {
		return entity.Customer{}, err
	}

	return customer, nil
}

func (c *Customer) generateJWTToken(cst entity.Customer) (string, error) {
	claims := &dto.Claims{
		ID:    cst.ID,
		Email: cst.Email,
		Name:  cst.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.cfg.JWTSecret))
}

func hashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		slog.Errorw("generate password failed", "stack_trace", err.Error())
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		slog.Errorw("compare password failed", "stack_trace", err.Error())
		return false
	}

	return true
}

package provider

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt interface {
	GenerateFromPassword(ctx context.Context, password []byte) ([]byte, error)
	CompareHashAndPassword(ctx context.Context, hashedPassword, password []byte) error
}

type BcryptImpl struct {
}

func NewBcrypt() Bcrypt {
	return &BcryptImpl{}
}

func (b *BcryptImpl) GenerateFromPassword(ctx context.Context, password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (b *BcryptImpl) CompareHashAndPassword(ctx context.Context, hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	Library "e-commerce/library"
)

var (
	ErrPkcs7stripDataEmpty           = errors.New("pkcs7: Data is empty")
	ErrPkcs7stripDataNotBlockAligned = errors.New("pkcs7: Data is not block-aligned")
	ErrPkcs7stripInvalidPadding      = errors.New("pkcs7: Invalid padding")
	ErrPkcs7padInvalidBlockSize      = errors.New("pkcs7: Invalid block size")
	ErrNoContentDecryption           = errors.New("decrypt: No Content")
	ErrInvalidDecryption             = errors.New("decrypt: Invalid Content")
)

type CustomCrypto interface {
	AesEncryptGcmNoPadding(strIn string, key string) (string, error)
	AesDecryptGcmNoPadding(aesKey string, encryptedBase64 string) (string, error)
	GenerateCEK(privateKeyBase64 string, publicKeyBase64 string) (string, error)
}
type CustomCryptoImpl struct {
	library Library.Library
}

func NewCustomCrypto(
	config Config.Config,
	library Library.Library,
) CustomCrypto {
	return &CustomCryptoImpl{
		library: library,
	}
}

// this function is just for testing ECDH (elliptic curve diffie-hellman)
func (c *CustomCryptoImpl) GenerateCEK(privateKeyBase64 string, publicKeyBase64 string) (string, error) {
	// Decode private key from base64
	privKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return Constants.NilString, fmt.Errorf(Constants.ErrDecodeBase64PrivateKey.Error(), err)
	}

	privKey, err := x509.ParseECPrivateKey(privKeyBytes)
	if err != nil {
		return Constants.NilString, fmt.Errorf(Constants.ErrParseECPrivateKey.Error(), err)
	}

	// Decode public key from base64
	pubKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return Constants.NilString, fmt.Errorf(Constants.ErrDecodeBase64PublicKey.Error(), err)
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return Constants.NilString, fmt.Errorf(Constants.ErrParsePublicKey.Error(), err)
	}

	ecPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return Constants.NilString, fmt.Errorf(Constants.ErrNotECDSAPublicKey.Error())
	}

	// Perform ECDH key exchange
	x, _ := privKey.Curve.ScalarMult(ecPubKey.X, ecPubKey.Y, privKey.D.Bytes())

	// Encode the shared secret as a base64 string
	sharedSecretBase64 := base64.StdEncoding.EncodeToString(x.Bytes())

	return sharedSecretBase64, nil
}

func (c *CustomCryptoImpl) AesEncryptGcmNoPadding(aesKey string, strIn string) (string, error) {
	// Decode AES key from base64
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return Constants.NilString, err
	}

	// Generate random IV (Initialization Vector)
	iv := make([]byte, 12) // IV length for AES-GCM is 12 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return Constants.NilString, err
	}

	// Create AES cipher block using key
	block, err := aes.NewCipher(key)
	if err != nil {
		return Constants.NilString, err
	}

	// Create GCM mode with AES block and IV
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return Constants.NilString, err
	}

	// Encrypt the plaintext
	ciphertext := aesGCM.Seal(nil, iv, []byte(strIn), nil)

	// Combine IV and ciphertext
	encryptedData := append(iv, ciphertext...)

	// Encode encrypted data to base64
	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedData)

	return encryptedBase64, nil
}

func (c *CustomCryptoImpl) AesDecryptGcmNoPadding(aesKey string, encryptedBase64 string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return Constants.NilString, err
	}

	// Decode encrypted data from base64
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return Constants.NilString, err
	}

	// Extract IV from encrypted data (first 12 bytes)
	iv := encryptedData[:12]
	ciphertext := encryptedData[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return Constants.NilString, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return Constants.NilString, err
	}

	// Decrypt the ciphertext using AES-GCM
	plaintext, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return Constants.NilString, err
	}

	return string(plaintext), nil
}

package constants

import (
	"errors"
)

var (
	ErrServeFailed            = errors.New("failed to run the apps")
	ErrConfiguration          = errors.New("failed to config")
	ErrFailedJSONMarshal      = errors.New("failed to convert json")
	ErrNoConfiguration        = errors.New("config not found")
	ErrPanic                  = errors.New("defer")
	ErrValidationRequired     = errors.New("%s must be filled")
	ErrValidationMax          = errors.New("%s surpass max character")
	ErrValidationOneOF        = errors.New("%s not match")
	ErrValidationMin          = errors.New("%s min %s character")
	ErrDataTypeInvalid        = errors.New("data type %s not fit")
	ErrAlphaNumeric           = errors.New("value not alphanumeric")
	ErrEmail                  = errors.New("value not email")
	ErrMultipleLogin          = errors.New("user still login, please logout first")
	ErrSomethingWentWrong     = errors.New("something went wrong")
	ErrFailedRemove           = errors.New("failed delete data")
	ErrJWTInvalidMethod       = errors.New("signing method not match")
	ErrInvalidToken           = errors.New("token not match")
	ErrDecodeAPIKey           = errors.New("failed to decode API key: %w")
	ErrFailedGenerateJWK      = errors.New("failed to create JWK: %w")
	ErrFailedEncryptJWT       = errors.New("failed to Encrypt JWT: %w")
	ErrFailedDecryptJWE       = errors.New("failed to decrypt JWE: %w")
	ErrUnmarshalClaim         = errors.New("failed to unmarshal claims: %w")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrConnectionPostgres     = errors.New("failed to open connection postgresql")
	ErrDuplicatedKey          = errors.New("duplicate data")
	ErrDuplicatedUser         = errors.New("user already exists")
	ErrOneStore               = errors.New("user already has store")
	ErrInactiveStore          = errors.New("user already has store, but inactive")
	ErrPrivateKeyIsEmpty      = errors.New("private key empty")
	ErrFailedSigning          = errors.New("signing Error")
	ErrInternalServerError    = errors.New("internal Server Error")
	ErrInvalidJWE             = errors.New("invalid JWE")
	ErrAuthorizationBearer    = errors.New("missing authorization bearer")
	ErrRoleNotMatch           = errors.New("role not match")
	ErrUserNotFound           = errors.New("user not found")
	ErrStoreNotFound          = errors.New("store not found")
	ErrWrongPassword          = errors.New("wrong password")
	ErrWrongEmailOrUsername   = errors.New("username or email is wrong")
	ErrConfigurationRedis     = errors.New("failed configuration redis")
	ErrValidation             = errors.New("validation error")
	ErrEmptyTraceID           = errors.New("empty trace id")
	ErrEmptyCredential        = errors.New("empty credential")
	ErrDecodeBase64PrivateKey = errors.New("failed to decode base64 private key: %v")
	ErrDecodeBase64PublicKey  = errors.New("failed to decode base64 public key: %v")
	ErrParseECPrivateKey      = errors.New("failed to parse EC private key: %v")
	ErrParsePublicKey         = errors.New("failed to parse public key: %v")
	ErrNotECDSAPublicKey      = errors.New("public key is not of type *ecdsa.PublicKey")
	ErrPerformRequest         = errors.New("failed to perform request: %v")
	ErrReadResponseBody       = errors.New("failed to read response body: %v")
	ErrPublicKeyIsEmpty       = errors.New("public key empty")
	ErrGenerateCEK            = errors.New("failed create content encryption key")
	ErrCEKIsEmpty             = errors.New("CEK is empty")
	ErrDuplicatedProxyNumber  = errors.New("duplicate proxy number")
	ErrCredentialNotFound     = errors.New("credential not found")
	ErrCardNotFound           = errors.New("card not found")
	ErrDesc                   = errors.New("error: %v")
	ErrFailedEncryptResponse  = errors.New("failed to encrypt response")
	ErrDecrypt                = errors.New("gagal melakukan decryption")
	ErrTransactionNotFound    = errors.New("transaction not found")
	ErrFailedGenerateTraceID  = errors.New("failed generate traceID")
)

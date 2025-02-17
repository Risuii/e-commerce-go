package usecase

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/authentication/delivery/dto"
	AuthenticationEntity "e-commerce/internal/authentication/domain/entity"
	AuthenticationRepository "e-commerce/internal/authentication/domain/repository"
	UserEntity "e-commerce/internal/user/domain/entity"
	UserRepository "e-commerce/internal/user/domain/repository"
	Library "e-commerce/library"
	BcryptPackage "e-commerce/pkg/bcrypt"
	CryptoPackage "e-commerce/pkg/crypto"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	ExecutionResultPackage "e-commerce/pkg/execution_result"
	JWEPackage "e-commerce/pkg/jwe"
	UtilsPackage "e-commerce/pkg/utils"
)

type LoginUsecase interface {
	Index(param AuthDTO.LoginParam) (string, error)
}

type LoginUsecaseImpl struct {
	crypto                   CryptoPackage.CustomCrypto
	bcrypty                  BcryptPackage.Bcrypt
	library                  Library.Library
	config                   Config.Config
	jwe                      JWEPackage.JWE
	userRepository           UserRepository.UserRepository
	authenticationRepository AuthenticationRepository.AuthenticationRepository
}

func NewLoginUsecase(
	jwe JWEPackage.JWE,
	bcrypty BcryptPackage.Bcrypt,
	crypto CryptoPackage.CustomCrypto,
	library Library.Library,
	config Config.Config,
	userRepository UserRepository.UserRepository,
	authenticationRepository AuthenticationRepository.AuthenticationRepository,
) LoginUsecase {
	return &LoginUsecaseImpl{
		jwe:                      jwe,
		bcrypty:                  bcrypty,
		crypto:                   crypto,
		library:                  library,
		config:                   config,
		userRepository:           userRepository,
		authenticationRepository: authenticationRepository,
	}
}

func (u *LoginUsecaseImpl) Index(param AuthDTO.LoginParam) (string, error) {
	path := "LoginUsecase:Index"

	// CHECKING REDIS
	oldToken, err := u.GetTokenInRedis(UtilsPackage.TernaryOperator(param.Username != Constants.NilString, param.Username, param.Email))
	if err != nil {
		return Constants.NilString, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// CHECKING TOKEN IS EXIST OR NOT
	if oldToken != nil {
		err := CustomErrorPackage.New(Constants.ErrMultipleLogin, Constants.ErrMultipleLogin, path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusUnauthorized)
		return Constants.NilString, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// MAKE CHANNEL
	UserChannel := make(chan ExecutionResultPackage.ExecutionResult)
	JWEChannel := make(chan ExecutionResultPackage.ExecutionResult)

	// DO GOROUTINE
	go u.GetDetailUsers(param.Username, param.Email, UserChannel)
	go u.GenerateJWE(param.Username, param.Email, JWEChannel)

	// GET USER DATA FROM CHANNEL
	UserDetail := <-UserChannel
	if err := UserDetail.GetError(); err != nil {
		return Constants.NilString, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	user := UserDetail.GetData().(*UserEntity.User)
	if user == nil {
		err := CustomErrorPackage.New(Constants.ErrUserNotFound, Constants.ErrUserNotFound, path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusUnauthorized)
		return Constants.NilString, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// GET JWE FROM CHANNEL
	JWEDetail := <-JWEChannel
	if err := JWEDetail.GetError(); err != nil {
		return Constants.NilString, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	token := JWEDetail.GetData().(string)

	// CHECK PASSWORD
	if err := u.CompareHash(user.Password, param.Password); err != nil {
		err := CustomErrorPackage.New(Constants.ErrWrongPassword, Constants.ErrWrongPassword, path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusUnauthorized)
		return Constants.NilString, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// INIT CHANNEL FOR STORE TOKEN TO REDIS
	storeJWEChannel := make(chan ExecutionResultPackage.ExecutionResult)

	// DO GOROUTINE
	go u.StoreToken(token, UtilsPackage.TernaryOperator(param.Username != Constants.NilString, param.Username, param.Email), u.config.GetConfig().JWE.ExpiryDuration, storeJWEChannel)

	return token, nil
}

func (u *LoginUsecaseImpl) EncryptEmail(email string) (string, error) {
	path := "LoginUsecase:EncryptEmail"

	// ENCRYPT EMAIL
	encryptedEmail, err := u.crypto.AesEncryptGcmNoPadding(u.config.GetConfig().EncryptKey.EncryptKey, email)
	if err != nil {
		return Constants.NilString, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// ENCODE ENCRYPTED EMAIL
	encodeEmail := base64.StdEncoding.EncodeToString([]byte(encryptedEmail))

	return encodeEmail, nil
}

func (u *LoginUsecaseImpl) GetDetailUsers(username, email string, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "LoginUsecase:GetDetailUsers"

	result := ExecutionResultPackage.ExecutionResult{}

	// CHECKING IF USERNAME AND EMAIL IS EXIST
	if username != Constants.NilString && email != Constants.NilString {
		email = Constants.NilString
	}

	// GET USER DETAIL
	entity, err := u.userRepository.GetDetailUsers(strings.ToLower(username), email)
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	// SET DATA USER TO CHANNEL
	result.SetResult(entity, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *LoginUsecaseImpl) GenerateJWE(username, email string, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "LoginUsecase:GenerateJWE"

	result := ExecutionResultPackage.ExecutionResult{}

	// INIT CLAIM
	claims := jwt.MapClaims{
		Constants.Username: username,
		Constants.Email:    email,
	}

	// GENERATE JWE
	token, err := u.jwe.JWEGenerateToken(claims, u.config.GetConfig().JWE.SecretKey)
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	// SET JWE TOKEN TO CHANNEL
	result.SetResult(token, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *LoginUsecaseImpl) CompareHash(hashedPassword, password string) error {
	path := "LoginUsecase:CompareHash"

	// DECODE PASSWORD
	decodedPassword, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		err := CustomErrorPackage.New(Constants.ErrInternalServerError, Constants.ErrInternalServerError, path, u.library)
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// DECRYPT PASSWORD
	decryptedPassword, err := u.crypto.AesDecryptGcmNoPadding(u.config.GetConfig().EncryptKey.EncryptKey, string(decodedPassword))
	if err != nil {
		err := CustomErrorPackage.New(Constants.ErrInternalServerError, Constants.ErrInternalServerError, path, u.library)
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// COMPARE HASH
	if err := u.bcrypty.CompareHashAndPassword(context.Background(), []byte(decryptedPassword), []byte(password)); err != nil {
		err := CustomErrorPackage.New(Constants.ErrInternalServerError, Constants.ErrInternalServerError, path, u.library)
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (u *LoginUsecaseImpl) StoreToken(token, username string, expire time.Duration, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "LoginUsecase:StoreToken"

	result := ExecutionResultPackage.ExecutionResult{}

	// INIT ENTITY
	entity := AuthenticationEntity.JWEToken{
		Username:         username,
		TOKEN:            token,
		JWT_TOKEN_PREFIX: Constants.UsernamePrefix,
	}

	if err := u.authenticationRepository.SetJWEToken(&entity, expire); err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	result.SetResult(nil, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *LoginUsecaseImpl) GetTokenInRedis(username string) (*AuthenticationEntity.JWEToken, error) {
	path := "LoginUsecase:GetTokenInRedis"

	token, err := u.authenticationRepository.GetJWETokenByKey(Constants.UsernamePrefix + username)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return token, nil
}

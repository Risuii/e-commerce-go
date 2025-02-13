package usecase

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/auth/delivery/dto"
	UserEntity "e-commerce/internal/auth/domain/entity"
	UserRepository "e-commerce/internal/auth/domain/repository"
	LogRepository "e-commerce/internal/logging/domain/repository"
	Library "e-commerce/library"
	BcryptPackage "e-commerce/pkg/bcrypt"
	CryptoPackage "e-commerce/pkg/crypto"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	ExecutionResultPackage "e-commerce/pkg/execution_result"
	RequestPackage "e-commerce/pkg/request_information"
)

type RegisterUseCase interface {
	Index(requestInfo RequestPackage.RequestInformation, param AuthDTO.RegisterParam, traceID string) error
}

type RegisterUseCaseImpl struct {
	crypto         CryptoPackage.CustomCrypto
	bcrypty        BcryptPackage.Bcrypt
	library        Library.Library
	config         Config.Config
	logRepository  LogRepository.LogActivity
	userRepository UserRepository.UserRepository
}

func NewRegisterUseCase(
	library Library.Library,
	crypto CryptoPackage.CustomCrypto,
	config Config.Config,
	bcrypty BcryptPackage.Bcrypt,
	logRepository LogRepository.LogActivity,
	userRepository UserRepository.UserRepository,
) RegisterUseCase {
	return &RegisterUseCaseImpl{
		crypto:         crypto,
		config:         config,
		library:        library,
		bcrypty:        bcrypty,
		logRepository:  logRepository,
		userRepository: userRepository,
	}
}

func (u *RegisterUseCaseImpl) Index(requestInfo RequestPackage.RequestInformation, param AuthDTO.RegisterParam, traceID string) error {
	path := "RegisterUsecase:Index"

	// MAKE CHANNEL
	IDChannel := make(chan ExecutionResultPackage.ExecutionResult)
	UserChannel := make(chan ExecutionResultPackage.ExecutionResult)
	PasswordChannel := make(chan ExecutionResultPackage.ExecutionResult)
	EmailChannel := make(chan ExecutionResultPackage.ExecutionResult)

	// DO GOROUTINE
	go u.GenerateID(IDChannel)
	go u.FindByUsername(param.Username, UserChannel)
	go u.HashPassword(param.Password, PasswordChannel)
	go u.EncryptEmail(param.Email, EmailChannel)

	// GET ID FROM CHANNEL
	IDDetail := <-IDChannel
	if err := IDDetail.GetError(); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	ID := IDDetail.GetData().(string)

	// GET USER DATA FROM CHANNEL
	UserDetail := <-UserChannel
	if err := UserDetail.GetError(); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	user := UserDetail.GetData().(*UserEntity.User)
	if user != nil {
		err := CustomErrorPackage.New(Constants.ErrDuplicatedUser, Constants.ErrDuplicatedUser, path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusConflict)
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// GET PASSWORD FROM CHANNEL
	PasswordDetail := <-PasswordChannel
	if err := PasswordDetail.GetError(); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	password := PasswordDetail.GetData().(string)

	// GET EMAIL FROM CHANNEL
	EmailDetail := <-EmailChannel
	if err := EmailDetail.GetError(); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	email := EmailDetail.GetData().(string)

	// INSERT DATA TO DB
	if err := u.Insert(ID, password, email, param); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (u *RegisterUseCaseImpl) GenerateID(resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "RegisterUsecase:GenerateID"

	result := ExecutionResultPackage.ExecutionResult{}

	id, err := u.library.GenerateUUID()
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	result.SetResult(id, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *RegisterUseCaseImpl) FindByUsername(username string, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "RegisterUsecase:FindByUsername"

	result := ExecutionResultPackage.ExecutionResult{}

	entity, err := u.userRepository.FindByUsername(strings.ToLower(username))
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	result.SetResult(entity, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *RegisterUseCaseImpl) HashPassword(password string, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "RegisterUsecase:HashPassword"

	result := ExecutionResultPackage.ExecutionResult{}

	hashedPassword, err := u.bcrypty.GenerateFromPassword(context.Background(), []byte(password))
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	encryptedPassword, err := u.crypto.AesEncryptGcmNoPadding(u.config.GetConfig().EncryptKey.EncryptKey, string(hashedPassword))
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	encodePassword := base64.StdEncoding.EncodeToString([]byte(encryptedPassword))

	result.SetResult(encodePassword, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *RegisterUseCaseImpl) EncryptEmail(email string, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "RegisterUsecase:EncryptEmail"

	result := ExecutionResultPackage.ExecutionResult{}

	encryptedEmail, err := u.crypto.AesEncryptGcmNoPadding(u.config.GetConfig().EncryptKey.EncryptKey, email)
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	encodeEmail := base64.StdEncoding.EncodeToString([]byte(encryptedEmail))

	result.SetResult(encodeEmail, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *RegisterUseCaseImpl) Insert(id, password, email string, param AuthDTO.RegisterParam) error {
	path := "RegisterUsecase:Insert"

	userData := UserEntity.User{
		Uuid:      id,
		Email:     email,
		Username:  strings.ToLower(param.Username),
		Password:  password,
		CreatedAt: time.Now().Format(Constants.YYYMMDDHHMMSS),
	}

	err := u.userRepository.Insert(&userData)
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

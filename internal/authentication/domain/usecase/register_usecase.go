package usecase

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/authentication/delivery/dto"
	LogRepository "e-commerce/internal/logging/domain/repository"
	UserEntity "e-commerce/internal/user/domain/entity"
	UserRepository "e-commerce/internal/user/domain/repository"
	Library "e-commerce/library"
	BcryptPackage "e-commerce/pkg/bcrypt"
	CryptoPackage "e-commerce/pkg/crypto"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	ExecutionResultPackage "e-commerce/pkg/execution_result"
)

type RegisterUseCase interface {
	Index(param AuthDTO.RegisterParam) error
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

func (u *RegisterUseCaseImpl) Index(param AuthDTO.RegisterParam) error {
	path := "RegisterUsecase:Index"

	// MAKE CHANNEL
	IDChannel := make(chan ExecutionResultPackage.ExecutionResult)
	UserChannel := make(chan ExecutionResultPackage.ExecutionResult)
	PasswordChannel := make(chan ExecutionResultPackage.ExecutionResult)

	// DO GOROUTINE
	go u.GenerateID(IDChannel)
	go u.GetDetailUsers(param.Username, UserChannel)
	go u.HashPassword(param.Password, PasswordChannel)

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

	// INSERT DATA TO DB
	if err := u.Insert(ID, password, param.Email, param); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (u *RegisterUseCaseImpl) GenerateID(resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "RegisterUsecase:GenerateID"

	// INIT CHANNEL
	result := ExecutionResultPackage.ExecutionResult{}

	// GENERATE UUID FOR ID
	id, err := u.library.GenerateUUID()
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	// SET UUID TO CHANNEL
	result.SetResult(id, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *RegisterUseCaseImpl) GetDetailUsers(username string, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "RegisterUsecase:GetDetailUsers"

	// INIT CHANNEL
	result := ExecutionResultPackage.ExecutionResult{}

	// GET USER DATA
	entity, err := u.userRepository.GetDetailUsers(strings.ToLower(username), Constants.NilString)
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	// SET USER DATA TO CHANNEL
	result.SetResult(entity, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *RegisterUseCaseImpl) HashPassword(password string, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "RegisterUsecase:HashPassword"

	// INIT CHANNEL
	result := ExecutionResultPackage.ExecutionResult{}

	// HASHING PASSWORD
	hashedPassword, err := u.bcrypty.GenerateFromPassword(context.Background(), []byte(password))
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	// ENCRYPT HASHED PASSWORD
	encryptedPassword, err := u.crypto.AesEncryptGcmNoPadding(u.config.GetConfig().EncryptKey.EncryptKey, string(hashedPassword))
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	// ENCODED ENCRYPTED HASHED PASSWORD
	encodePassword := base64.StdEncoding.EncodeToString([]byte(encryptedPassword))

	// SET PASSWORD TO CHANNEL
	result.SetResult(encodePassword, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *RegisterUseCaseImpl) Insert(id, password, email string, param AuthDTO.RegisterParam) error {
	path := "RegisterUsecase:Insert"

	// INIT ENTITY USER
	userData := UserEntity.User{
		Uuid:      id,
		Email:     email,
		Username:  strings.ToLower(param.Username),
		Password:  password,
		CreatedAt: time.Now().Format(Constants.YYYMMDDHHMMSS),
	}

	// INSERT ENTITY USER TO DB
	err := u.userRepository.Insert(&userData)
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

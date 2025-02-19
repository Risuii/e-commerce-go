package source

import (
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"

	Constants "e-commerce/constants"
	UserModel "e-commerce/internal/user/data/model"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	EcommercePackage "e-commerce/pkg/data_sources/e-commerce"
)

type UserPersistent interface {
	GetDetail(username, email string) (*UserModel.User, error)
	Create(param *UserModel.User) error
	UpdateLastLogin(userID, lastLogin string) error
}

type UserImpl struct {
	library     Library.Library
	dbEcommerce EcommercePackage.Ecommerce
}

func NewUserPersistent(
	library Library.Library,
	dbEcommerce EcommercePackage.Ecommerce,
) UserPersistent {
	return &UserImpl{
		library:     library,
		dbEcommerce: dbEcommerce,
	}
}

func (s *UserImpl) GetDetail(username, email string) (*UserModel.User, error) {
	path := "UserPersistent:GetDetail"

	query := `
		SELECT 
			uuid, 
			email, 
			email_verified, 
			username, 
			password, 
			last_login, 
			created_at, 
			updated_at, 
			deleted_at 
		FROM 
			users 
	`

	// INIT QUERY WHERE
	queryWheres := []string{}

	// ADD VARIABLE QUERY ARGS
	queryArgs := []any{}

	if username != Constants.NilString {
		queryWheres = append(queryWheres, "username = ?")
		queryArgs = append(queryArgs, username)
	}

	if email != Constants.NilString {
		queryWheres = append(queryWheres, "email = ?")
		queryArgs = append(queryArgs, email)
	}

	if len(queryWheres) > 0 {
		query += " WHERE " + s.library.StringsJoin(queryWheres, " OR ")
	}

	var model *UserModel.User
	err := s.dbEcommerce.GetConnection().Raw(query, queryArgs...).Scan(&model).Error
	if err == nil {
		return model, nil
	}

	// THREAT "ErrRecordNotFound" AS SUCCESS WITH NULL DATA
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return model, CustomErrorPackage.New(Constants.ErrSomethingWentWrong, err, path, s.library)
}

func (s *UserImpl) Create(param *UserModel.User) error {
	path := "UserPersistent:Create"

	err := s.dbEcommerce.GetConnection().Exec(`
		INSERT INTO users (
			uuid,
			email,
			email_verified,
			username,
			password,
			created_at
		) 
		VALUES (
			?, ?, ?, ?, ?, ?
		)
	`,
		param.Uuid,
		param.Email,
		false,
		param.Username,
		param.Password,
		param.CreatedAt,
	).Error

	if err == nil {
		return nil
	}

	//ASSERT POSGRES ERROR
	var postgresError *pq.Error
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return CustomErrorPackage.New(Constants.ErrDuplicatedKey, err, path, s.library)
	}

	return CustomErrorPackage.New(Constants.ErrSomethingWentWrong, err, path, s.library)
}

func (s *UserImpl) UpdateLastLogin(userID, lastLogin string) error {
	path := "UserPersistent:UpdateLastLogin"

	err := s.dbEcommerce.GetConnection().Exec(`
		UPDATE
			users
		SET
			last_login = ?
		WHERE
			uuid = ?
	`,
		lastLogin,
		userID,
	).Error

	if err == nil {
		return nil
	}
	//ASSERT POSGRES ERROR
	var postgresError *pq.Error
	if errors.As(err, &postgresError) && postgresError.Code == Constants.PostgresErrorCode {
		return CustomErrorPackage.New(Constants.ErrDuplicatedKey, err, path, s.library)
	}

	return CustomErrorPackage.New(Constants.ErrSomethingWentWrong, err, path, s.library)
}

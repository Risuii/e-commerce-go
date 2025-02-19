package source

import (
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"

	Constants "e-commerce/constants"
	StoreModel "e-commerce/internal/store/data/model"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	EcommercePackage "e-commerce/pkg/data_sources/e-commerce"
)

type StorePersistent interface {
	Create(param *StoreModel.Store) error
	Update(param *StoreModel.Store) error
	GetStore(userID string) (*StoreModel.Store, error)
}

type StorePersistentImpl struct {
	library     Library.Library
	dbEcommerce EcommercePackage.Ecommerce
}

func NewStorePersistent(
	library Library.Library,
	dbEcommerce EcommercePackage.Ecommerce,
) StorePersistent {
	return &StorePersistentImpl{
		library:     library,
		dbEcommerce: dbEcommerce,
	}
}

func (s *StorePersistentImpl) Create(param *StoreModel.Store) error {
	path := "StorePersistent:Create"

	err := s.dbEcommerce.GetConnection().Exec(`
		INSERT INTO stores (
			store_id,
			store_name,
			description,
			user_id,
			status,
			created_at
		) 
		VALUES (
			?, ?, ?, ?, ?, ?
		)
	`,
		param.StoreID,
		param.StoreName,
		param.Description,
		param.UserID,
		param.Status,
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

func (s *StorePersistentImpl) Update(param *StoreModel.Store) error {
	path := "StorePersistent:Update"

	err := s.dbEcommerce.GetConnection().Exec(`
		UPDATE
			stores
		SET
			store_name = ?,
			description = ?,
			status = ?,
			updated_at = ?
		WHERE
			store_id = ?
	`,
		param.StoreName,
		param.Description,
		param.Status,
		param.UpdatedAt,
		param.StoreID,
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

func (s *StorePersistentImpl) GetStore(userID string) (*StoreModel.Store, error) {
	path := "StorePersistent:GetStore"

	query := `
		SELECT
			store_id,
			store_name,
			description,
			user_id,
			status,
			created_at,
			updated_at
		FROM
			stores
	`

	// INIT QUERY WHERE
	queryWheres := []string{}

	// ADD VARIABLE QUERY ARGS
	queryArgs := []any{}

	if userID != Constants.NilString {
		queryWheres = append(queryWheres, "user_id = ?")
		queryArgs = append(queryArgs, userID)
	}

	if len(queryWheres) > 0 {
		query += " WHERE " + s.library.StringsJoin(queryWheres, " AND ")
	}

	var model *StoreModel.Store
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

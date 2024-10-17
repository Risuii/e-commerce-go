package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	Constants "e-commerce/constants"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
)

func GetCustomErrorMessage(validationErrors validator.ValidationErrors, callback func(anyParam interface{}) interface{}, fields ...string) []map[string]interface{} {
	var errors []map[string]interface{}
	for i, verr := range validationErrors {
		fieldName := verr.Field()
		if len(fields) == len(validationErrors) {
			fieldName = fields[i]
		}

		switch verr.Tag() {
		case Constants.ValidationRequired:
			errors = append(errors, map[string]interface{}{
				Constants.Field:   verr.Field(),
				Constants.Message: fmt.Sprintf(Constants.ErrValidationRequired.Error(), fieldName),
			})
			break
		case Constants.ValidationGt:
			errors = append(errors, map[string]interface{}{
				Constants.Field:   verr.Field(),
				Constants.Message: fmt.Sprintf(Constants.ErrValidationRequired.Error(), fieldName),
			})
			break
		case Constants.ValidationMax:
			errors = append(errors, map[string]interface{}{
				Constants.Field:   verr.Field(),
				Constants.Message: fmt.Sprintf(Constants.ErrValidationMax.Error(), fieldName),
			})
			break
		case Constants.ValidationOneOf:
			errors = append(errors, map[string]interface{}{
				Constants.Field:   verr.Field(),
				Constants.Message: fmt.Sprintf(Constants.ErrValidationOneOF.Error(), fieldName),
			})
			break
		case Constants.ValidationRequiredWithout:
			errors = append(errors, map[string]interface{}{
				Constants.Field:   verr.Field(),
				Constants.Message: fmt.Sprintf(Constants.ErrValidationRequired.Error(), fieldName),
			})
			break
		default:
			errors = append(errors, map[string]interface{}{
				Constants.Field:   verr.Field(),
				Constants.Message: Constants.ErrSomethingWentWrong.Error(),
			})
			break
		}
	}

	return errors
}

func GetCustomFieldName(validationErrors validator.ValidationErrors, tag string, getField func(field string) *reflect.StructField) *[]string {
	var fields []string
	for _, v := range validationErrors {
		field := getField(v.Field())
		if field == nil {
			fields = append(fields, v.Field())
			continue
		}

		fieldName, ok := (*field).Tag.Lookup(tag)
		if !ok {
			fields = append(fields, v.Field())
			continue
		}

		fields = append(fields, fieldName)
	}

	return &fields
}

func TernaryOperator[T interface{}](comparator bool, trueCondition T, falseCondition T) T {
	if comparator {
		return trueCondition
	}

	return falseCondition
}

func TernaryOperatorPromise[T interface{}](comparator bool, trueCallback func() T, falseCallback func() T) T {
	if comparator {
		return trueCallback()
	}

	return falseCallback()
}

// CONVERT DATETIME TO STRING
func DateTimeToString(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02 15:04:05")
}

// SPLIT LATITUDE AND LONGITUDE FROM STRING DATA
func SplitLatLong(latLong string) (latitude string, longitude string) {
	if latLong == "" {
		return "", ""
	}
	latlongList := strings.Split(latLong, ", ")
	return latlongList[0], latlongList[1]
}

func CatchPanic(path string, library Library.Library) {
	if err := recover(); err != nil {
		err = CustomErrorPackage.New(Constants.ErrPanic, err.(error), path, library)
	}
}

func HourToYearAndMonth(hours int) (year int, month int) {
	convMonth := hours / 24 / 30
	year = convMonth / 12
	month = TernaryOperator(year > 0, (convMonth - (year * 12)), month)
	return year, month
}

// Format time month or day to bahasa
func TimeFormatID(date string) string {
	day := strings.NewReplacer(
		"Monday", "Senin",
		"Tuesday", "Selasa",
		"Wednesday", "Rabu",
		"Thursday", "Kamis",
		"Friday", "Jum'at",
		"Saturday", "Sabtu",
		"Sunday", "Minggu",
	)

	date = day.Replace(date)

	month := strings.NewReplacer(
		"January", "Januari",
		"February", "Februari",
		"March", "Maret",
		"April", "April",
		"May", "Mei",
		"June", "Juni",
		"July", "Juli",
		"August", "Agustus",
		"September", "September",
		"October", "Oktober",
		"November", "November",
		"December", "Desember",
	)

	date = month.Replace(date)

	return date
}

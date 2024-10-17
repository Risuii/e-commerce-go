package utils

import (
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"

	Constants "e-commerce/constants"
)

// Validate date string is a valid date
func IsValidDate(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("tanggal tidak valid")
	}
	return nil
}

// Validate time string is a valid time
func IsValidTime(times string) error {
	_, err := time.Parse("15:04", times)
	if err != nil {
		return fmt.Errorf("waktu tidak valid")
	}
	return nil
}

// Validate if date input is a weekend
func IsWeekend(date time.Time) bool {
	if date.Weekday() == 6 || date.Weekday() == 0 {
		return true
	}
	return false
}

func MinLength(fl validator.FieldLevel) bool {
	minTag := fl.Param() // GET VALUE FROM TAG MIN
	min, err := strconv.Atoi(minTag)
	if err != nil {
		// Jika terjadi kesalahan parsing, kembalikan false
		return false
	}

	return len(fl.Field().String()) >= min
}

// HasDigit checks if the param contains at least one digit.
func HasDigit(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	paramDecode, err := base64.StdEncoding.DecodeString(param)
	if err != nil {
		return false
	}

	for _, char := range string(paramDecode) {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

// HasUpper checks if the param contains at least one uppercase letter.
func HasUpper(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	for _, char := range param {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

// HasLower checks if the param contains at least one lowercase letter.
func HasLower(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	for _, char := range param {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

// HasSpecialChars checks if the param contains at least one special character.
func HasSpecialChars(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	paramEncode, _ := base64.StdEncoding.DecodeString(param)
	specialChars := "!@#$%^&*?_~-£().," // Ganti string dengan karakter khusus yang diinginkan
	for _, char := range string(paramEncode) {
		if strings.ContainsRune(specialChars, char) {
			return true
		}
	}
	return false
}

func YYYYMMDDGTEField(fl validator.FieldLevel) bool {
	value, _, _, ok := fl.GetStructFieldOK2()

	if !ok {
		panic(Constants.ErrCustomValidator)
	}

	field := fl.Field().Interface().(string)

	target := value.Interface().(string)

	if target == "" {
		panic(Constants.ErrCustomValidator)
	}

	fieldTime, err := time.Parse(Constants.YYYMMDD, field)
	if err != nil {
		panic(Constants.ErrCustomValidator)
	}

	fieldValueTime, err := time.Parse(Constants.YYYMMDD, target)
	if err != nil {
		panic(Constants.ErrCustomValidator)
	}

	return fieldTime.Equal(fieldValueTime) || fieldTime.After(fieldValueTime)
}

func ValidateFileSize(fl validator.FieldLevel) bool {
	field := fl.Field().Interface().(multipart.FileHeader)
	param := fl.Param()
	max, err := strconv.Atoi(param)
	if err != nil {
		panic(Constants.ErrCustomValidator)
	}
	// INIT MAX SIZE
	maxFileSize := int64(max * 1024 * 1024)
	return field.Size <= maxFileSize
}

func ValidateMultipleFileSize(fl validator.FieldLevel) bool {
	fields := fl.Field().Interface().([]multipart.FileHeader)
	for _, field := range fields {
		param := fl.Param()
		max, err := strconv.Atoi(param)
		if err != nil {
			panic(Constants.ErrCustomValidator)
		}
		// INIT MAX SIZE
		maxFileSize := int64(max * 1024 * 1024)
		if field.Size > maxFileSize {
			return false
		}
	}

	return true
}

func ValidateMimeType(fl validator.FieldLevel) bool {
	field := fl.Field().Interface().(multipart.FileHeader)
	param := fl.Param()

	return strings.ToLower(field.Header.Get("Content-Type")) == param
}

func ValidateMultipleMimeType(fl validator.FieldLevel) bool {
	fields := fl.Field().Interface().([]multipart.FileHeader)
	for _, field := range fields {
		param := fl.Param()
		// INIT MAX SIZE
		if strings.ToLower(field.Header.Get("Content-Type")) != param {
			return false
		}
	}

	return true
}

func ValidateDateLessThan(fl validator.FieldLevel) bool {
	dateString := fl.Field().String()

	// Parsing string tanggal ke objek time.Time
	date, err := time.Parse(Constants.YYYMMDD, dateString)
	if err != nil {
		return false
	}

	// Validasi tanggal tidak boleh kurang dari tanggal sekarang
	return date.After(time.Now())
}

func ValidateDatetimeLessThan(fl validator.FieldLevel) bool {
	dateString := fl.Field().String()

	// Parsing string tanggal ke objek time.Time
	date, err := time.Parse(Constants.YYYMMDDHHMMSS, dateString)
	if err != nil {
		return false
	}

	// Validasi tanggal tidak boleh kurang dari tanggal sekarang
	return date.After(time.Now())
}

func ValidateZero(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}

	return value >= 0
}

func Date(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	// IF STRING CAN BE PARSED WITH DATE FORMAT THAN FIELD IS DATE
	_, err := time.Parse(Constants.YYYMMDD, param)
	return err == nil
}

func Time(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	// IF STRING CAN BE PARSED WITH TIME FORMAT THAN FIELD IS TIME
	_, err := time.Parse(Constants.HHMMSS, param)
	return err == nil
}

func DateGreaterOrEqualThanNow(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	date, err := time.Parse(Constants.YYYMMDD, param)
	if err != nil {
		panic(Constants.ErrCustomValidator)
	}

	strDateNow := time.Now().Format(Constants.YYYMMDD)
	dateNow, err := time.Parse(Constants.YYYMMDD, strDateNow)
	if err != nil {
		panic(Constants.ErrCustomValidator)
	}

	return date.Equal(dateNow) || date.After(dateNow)
}

func TimeGreaterOrEqualThanNow(fl validator.FieldLevel) bool {
	value, _, _, ok := fl.GetStructFieldOK2()

	if !ok {
		panic(Constants.ErrCustomValidator)
	}

	timeField := fl.Field().Interface().(string)

	dateField := value.Interface().(string)

	if dateField == "" {
		panic(Constants.ErrCustomValidator)
	}

	strDateTime := dateField + " " + timeField

	dateTime, err := time.Parse(Constants.YYYMMDDHHMMSS, strDateTime)
	if err != nil {
		// CASE IF INPUT DATE IS NOT DATE THAN RETURN TRUE (DON'T CHECK)
		return true
	}

	strDatetimeNow := time.Now().Format(Constants.YYYMMDDHHMMSS)
	datetimeNow, err := time.Parse(Constants.YYYMMDDHHMMSS, strDatetimeNow)
	if err != nil {
		panic(Constants.ErrCustomValidator)
	}

	return dateTime.Equal(datetimeNow) || dateTime.After(datetimeNow)
}

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	param := fl.Field().String()

	regex := `^(?:\+?62|0)8[0-9]{8,11}$`

	re := regexp.MustCompile(regex)

	return re.MatchString(param)
}

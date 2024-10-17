package library

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"golang.org/x/crypto/bcrypt"

	Sendinblue "github.com/sendinblue/APIv3-go-library/v2/lib"
)

type Type int
type Library interface {
	LoadEnv(filenames ...string) error
	Getenv(key string) string
	ReadEnv(cfg interface{}) error
	Getwd() (string, error)
	ReadConfig(path string, cfg interface{}) error
	GetNow() time.Time
	JsonMarshal(v interface{}) ([]byte, error)
	JsonUnmarshal(data []byte, v interface{}) error
	RandRead(b []byte) (n int, err error)
	AESNewCipher(key []byte) (cipher.Block, error)
	GetAES256CBCBlockSize() int
	Base64DecodeString(s string) ([]byte, error)
	Base64EncodeString(src []byte) string
	GetSlicesByteLen(v []byte) int
	ParseInt(v byte) int
	HasSuffix(s []byte, suffix []byte) bool
	JWTSignedString(jwtToken *jwt.Token, key interface{}) (string, error)
	CastJWTSigningMethodHMAC(jwtToken *jwt.Token) (*jwt.SigningMethodHMAC, bool)
	CastJWTMapClaims(jwtToken *jwt.Token) (jwt.MapClaims, bool)
	GetNewRotateFileHook(config rotatefilehook.RotateFileConfig) (logrus.Hook, error)
	SendTransacEmail(client *Sendinblue.APIClient, ctx context.Context, sendSmtpEmail Sendinblue.SendSmtpEmail) error
	GenerateUUID() (string, error)
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
	AtoI(param string) (int, error)
	FileInfo(directory string) (fs.FileInfo, error)
	IsNotExist(err error) bool
	CreateFile(filePath string) (*os.File, error)
	CreateDirectory(path string, perm fs.FileMode) error
	ReadFile(path string) ([]byte, error)
	CopyContentFile(dst io.Writer, src io.Reader) (written int64, err error)
	Contains(s string, substr string) bool
}

type LibraryImpl struct {
	loadEnv                func(filenames ...string) error
	getEnv                 func(key string) string
	readEnv                func(cfg interface{}) error
	getWd                  func() (string, error)
	readConfig             func(path string, cfg interface{}) error
	getNow                 func() time.Time
	jsonMarshal            func(v interface{}) ([]byte, error)
	jsonUnmarshal          func(data []byte, v interface{}) error
	randRead               func(b []byte) (n int, err error)
	aesNewCipher           func(key []byte) (cipher.Block, error)
	aes256CbcBlockSize     int
	base64DecodeString     func(s string) ([]byte, error)
	base64EncodeString     func(src []byte) string
	byteHasSuffix          func(s []byte, suffix []byte) bool
	rotateFileHook         func(config rotatefilehook.RotateFileConfig) (logrus.Hook, error)
	newRandomUUID          func() (uuid.UUID, error)
	generatePassword       func(password []byte, length int) ([]byte, error)
	compareHashAndPassword func(hashedPassword []byte, password []byte) error
	atoI                   func(param string) (int, error)
	fileInfo               func(directory string) (fs.FileInfo, error)
	isNotExist             func(err error) bool
	createFile             func(filePath string) (*os.File, error)
	mkdirAll               func(path string, perm fs.FileMode) error
	readFile               func(name string) ([]byte, error)
	copyContentFile        func(dst io.Writer, src io.Reader) (written int64, err error)
	contains               func(s string, substr string) bool
}

func New() Library {
	return &LibraryImpl{
		loadEnv:                godotenv.Load,
		getEnv:                 os.Getenv,
		readEnv:                cleanenv.ReadEnv,
		getWd:                  os.Getwd,
		readConfig:             cleanenv.ReadConfig,
		getNow:                 time.Now,
		jsonMarshal:            json.Marshal,
		jsonUnmarshal:          json.Unmarshal,
		randRead:               rand.Read,
		aesNewCipher:           aes.NewCipher,
		aes256CbcBlockSize:     16,
		base64DecodeString:     base64.StdEncoding.DecodeString,
		base64EncodeString:     base64.StdEncoding.EncodeToString,
		byteHasSuffix:          bytes.HasSuffix,
		rotateFileHook:         rotatefilehook.NewRotateFileHook,
		newRandomUUID:          uuid.NewRandom,
		generatePassword:       bcrypt.GenerateFromPassword,
		compareHashAndPassword: bcrypt.CompareHashAndPassword,
		atoI:                   strconv.Atoi,
		fileInfo:               os.Stat,
		isNotExist:             os.IsNotExist,
		createFile:             os.Create,
		mkdirAll:               os.MkdirAll,
		readFile:               os.ReadFile,
		copyContentFile:        io.Copy,
		contains:               strings.Contains,
	}
}

func (l *LibraryImpl) LoadEnv(filenames ...string) error {
	return l.loadEnv(filenames...)
}

func (l *LibraryImpl) Getenv(key string) string {
	return l.getEnv(key)
}

func (l *LibraryImpl) ReadEnv(cfg interface{}) error {
	return l.readEnv(cfg)
}

func (l *LibraryImpl) Getwd() (string, error) {
	return l.getWd()
}

func (l *LibraryImpl) ReadConfig(path string, cfg interface{}) error {
	return l.readConfig(path, cfg)
}

func (l *LibraryImpl) GetNow() time.Time {
	return l.getNow()
}

func (l *LibraryImpl) JsonMarshal(v interface{}) ([]byte, error) {
	return l.jsonMarshal(v)
}

func (l *LibraryImpl) JsonUnmarshal(data []byte, v interface{}) error {
	return l.jsonUnmarshal(data, v)
}

func (l *LibraryImpl) RandRead(b []byte) (n int, err error) {
	return l.randRead(b)
}

func (l *LibraryImpl) AESNewCipher(key []byte) (cipher.Block, error) {
	return l.aesNewCipher(key)
}

func (l *LibraryImpl) GetAES256CBCBlockSize() int {
	return l.aes256CbcBlockSize
}

func (l *LibraryImpl) Base64DecodeString(s string) ([]byte, error) {
	return l.base64DecodeString(s)
}

func (l *LibraryImpl) Base64EncodeString(s []byte) string {
	return l.base64EncodeString(s)
}

func (l *LibraryImpl) GetSlicesByteLen(v []byte) int {
	return len(v)
}

func (l *LibraryImpl) ParseInt(v byte) int {
	return int(v)
}

func (l *LibraryImpl) HasSuffix(s []byte, suffix []byte) bool {
	return l.byteHasSuffix(s, suffix)
}

func (l *LibraryImpl) JWTSignedString(jwtToken *jwt.Token, key interface{}) (string, error) {
	return jwtToken.SignedString(key)
}

func (l *LibraryImpl) CastJWTSigningMethodHMAC(jwtToken *jwt.Token) (*jwt.SigningMethodHMAC, bool) {
	method, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, ok
	}

	return method, ok
}

func (l *LibraryImpl) CastJWTMapClaims(jwtToken *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ok
	}

	return claims, ok
}

func (l *LibraryImpl) GetNewRotateFileHook(config rotatefilehook.RotateFileConfig) (logrus.Hook, error) {
	return l.rotateFileHook(config)
}

func (l *LibraryImpl) SendTransacEmail(client *Sendinblue.APIClient, ctx context.Context, sendSmtpEmail Sendinblue.SendSmtpEmail) error {
	_, _, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, sendSmtpEmail)
	return err
}

func (l *LibraryImpl) GenerateUUID() (string, error) {
	// GENERATE RANDOM UUID
	newUUID, err := l.newRandomUUID()
	return newUUID.String(), err
}

func (l *LibraryImpl) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return l.generatePassword(password, cost)
}

func (l *LibraryImpl) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return l.compareHashAndPassword(hashedPassword, password)
}

func (l *LibraryImpl) AtoI(param string) (int, error) {
	return l.atoI(param)
}

func (l *LibraryImpl) CreateFile(filePath string) (*os.File, error) {
	return l.createFile(filePath)
}

func (l *LibraryImpl) FileInfo(directory string) (fs.FileInfo, error) {
	return l.fileInfo(directory)
}

func (l *LibraryImpl) IsNotExist(err error) bool {
	return l.isNotExist(err)
}

func (l *LibraryImpl) CreateDirectory(path string, perm fs.FileMode) error {
	return l.mkdirAll(path, perm)
}

func (l *LibraryImpl) ReadFile(path string) ([]byte, error) {
	return l.readFile(path)
}

func (l *LibraryImpl) CopyContentFile(dst io.Writer, src io.Reader) (written int64, err error) {
	return l.copyContentFile(dst, src)
}

func (l *LibraryImpl) Contains(s string, substr string) bool {
	return l.contains(s, substr)
}

package environments

type EncryptKey struct {
	EncryptKey string `env-required:"true" env:"ENCRYPT_KEY"`
}

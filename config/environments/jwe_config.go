package environments

import "time"

type JWEConfig struct {
	SecretKey      string        `env-required:"true" env:"JWE_SECRET_KEY"`
	KID            string        `env-required:"true" env:"JWE_HEADER_VALUE_KID"`
	ExpiryDuration time.Duration `env-required:"true" env:"JWE_EXPIRY_DURATION"`
}

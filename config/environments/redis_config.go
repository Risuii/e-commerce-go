package environments

type RedisConfig struct {
	Host     string `env-required:"true" env:"REDIS_HOST"`
	Port     string `env-required:"true" env:"REDIS_PORT"`
	Password string `env-required:"true" env:"REDIS_PASSWORD"`
	Database string `env-required:"true" env:"REDIS_DATABASE"`
}

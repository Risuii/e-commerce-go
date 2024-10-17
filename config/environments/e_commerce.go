package environments

type Ecommerce struct {
	Host     string `env-required:"true" env:"DB_HOST"`
	Port     string `env-required:"true" env:"DB_PORT"`
	User     string `env-required:"true" env:"DB_USER"`
	Password string `env-required:"true" env:"DB_PASSWORD"`
	Database string `env-required:"true" env:"DB_DATABASE"`
}

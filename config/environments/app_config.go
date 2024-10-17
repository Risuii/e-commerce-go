package environments

type AppConfigModel struct {
	AppPort          string `env-required:"true" env:"APP_PORT"`
	LocationTimezone string `env-required:"true" env:"APP_LOCATION_TIMEZONE"`
}

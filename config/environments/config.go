package environments

type ConfigModel struct {
	App        AppConfigModel
	DB         DbConfigModel
	DBCon      Ecommerce
	EncryptKey EncryptKey
	JWE        JWEConfig
	Redis      RedisConfig
}

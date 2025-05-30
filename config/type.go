package config

type Config struct {
	App      AppConfig      `yaml:"app" validate:"required"`
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Redis    RedisConfig    `yaml:"redis" validate:"required"`
	Secret   SecretConfig   `yaml:"secret" validate:"required"`
	Product  ProductConfig  `yaml:"product" validate:"required"`
}

type AppConfig struct {
	Port string `yaml:"port" validate:"required"`
}

type ProductConfig struct {
	Host string `yaml:"host" validate:"required"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
}
type RedisConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}
type SecretConfig struct {
	JWTSecret string `yaml:"jwt_secret" validate:"required"`
}

package config

type Config struct {
	App         AppConfig             `mapstructure:"app" validate:"required"`
	Database    DatabaseConfig        `mapstructure:"database" validate:"required"`
	Redis       RedisConfig           `mapstructure:"redis" validate:"required"`
	KafkaConfig KafkaConfig           `mapstructure:"kafka" validate:"required"`
	PGAConfig   PaymentGateAwayConfig `mapstructure:"pga" validate:"required"`
}

type AppConfig struct {
	Port string `mapstructure:"port" validate:"required"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	Name     string `mapstructure:"name" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
}

type RedisConfig struct {
	Port     string `mapstructure:"port" validate:"required"`
	Host     string `mapstructure:"host" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}

type KafkaConfig struct {
	Broker      string            `mapstructure:"broker" validate:"required"`
	KafkaTopics map[string]string `mapstructure:"topics" validate:"required"`
}

type PaymentGateAwayConfig struct {
	ApiKey       string `mapstructure:"api_key" validate:"required"`
	WebhookToken string `mapstructure:"webhook_token" validate:"required"`
}

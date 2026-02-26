package configs

import (
	"log"

	"github.com/spf13/viper"
)

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`
}

type SqlDb struct {
	Host           string `json:"host"`
	Port           string `json:"port"`
	User           string `json:"user"`
	Password       string `json:"password"`
	DbName         string `json:"db_name"`
	DbMaxOpenCons  string `json:"db_max_open_cons"`
	DbIdleOpenCons string `json:"db_idle_open_cons"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Supabase struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

type Config struct {
	App      App      `json:"app"`
	SqlDb    SqlDb    `json:"sql_db"`
	Redis    Redis    `json:"redis"`
	RabbitMQ RabbitMQ `json:"rabbitmq"`
	Supabase Supabase `json:"supabase"`
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"),
		},
		SqlDb: SqlDb{
			Host:           viper.GetString("DATABASE_HOST"),
			Port:           viper.GetString("DATABASE_PORT"),
			User:           viper.GetString("DATABASE_USER"),
			Password:       viper.GetString("DATABASE_PASSWORD"),
			DbName:         viper.GetString("DATABASE_NAME"),
			DbMaxOpenCons:  viper.GetString("DATABASE_MAX_OPEN_CONNECTION"),
			DbIdleOpenCons: viper.GetString("DATABASE_MAX_IDLE_CONNECTION"),
		},
		Redis: Redis{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
		RabbitMQ: RabbitMQ{
			Host:     viper.GetString("RABBITMQ_HOST"),
			Port:     viper.GetString("RABBITMQ_PORT"),
			Username: viper.GetString("RABBITMQ_USER"),
			Password: viper.GetString("RABBITMQ_PASS"),
		},
		Supabase: Supabase{
			Url: viper.GetString("SUPABASE_URL"),
			Key: viper.GetString("SUPABASE_KEY"),
		},
	}
}

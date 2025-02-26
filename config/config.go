package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Configs struct {
	Postgres    PostgresConfig  `json:"postgres"`
	Elastic     ElasticConfig   `json:"elastic"`
	Cassandra   CassandraConfig `json:"cassandra"`
	Rabbit      RabbitConfig    `json:"rabbit"`
	Env         Env             `json:"env"`
	LogrusLevel uint8           `json:"logrus_level"`

	BIApiKey string `json:"bi_api_key"`
}

type Env struct {
	Mode                 string
	Namespace            string
	ForteKeysServiceName string `json:"forte_keys_service_name"`
	JwtSecret            string `json:"jwt_secret"`
	AccessTokenExpMin    int    `json:"access_token_exp_min"`
	RefreshTokenExpMin   int    `json:"refresh_token_exp_min"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `env:"user"`
	Password string `env:"password"`
	Database string `env:"database"`
	Secret   string `env:"secret"`
}

type ElasticConfig struct {
	ConnectionURL []string `json:"connection_url"`
	Login         string   `json:"login"`
	Password      string   `json:"password"`
}

type CassandraConfig struct {
	ConnectionIP []string `json:"connection_ip"`
}

type RabbitConfig struct {
	Host        string `json:"host"`
	VirtualHost string `json:"virtual_host"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	LogLevel    uint8  `json:"log_level"`
}

type DarLogisticsConfig struct {
	DeliveryPriceURL  string `json:"delivery_price_url"`
	ToDeliveryURL     string `json:"to_delivery_url"`
	DeliveryStatusURL string `json:"delivery_status_url"`
}

var AllConfigs *Configs

func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден, используем переменные среды")
	}

	accessTokenExpMin, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP_MIN"))
	refreshTokenExpMin, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP_MIN"))

	AllConfigs = &Configs{
		Env: Env{
			JwtSecret:          os.Getenv("JWT_SECRET"),
			AccessTokenExpMin:  accessTokenExpMin,
			RefreshTokenExpMin: refreshTokenExpMin,
		},
	}

	log.Println("Конфигурация загружена:", AllConfigs)
	return nil
}

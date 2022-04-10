package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	RabbitMQ RabbitMQ
}

type RabbitMQ struct {
	URL            string
	QueryName      string
	CountConsumers int
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		return nil, err
	}
	err = validator.New().Struct(c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

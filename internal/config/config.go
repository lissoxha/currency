package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string        `mapstructure:"PORT"`
	WeatherAPIKey  string        `mapstructure:"WEATHER_API_KEY"`
	WeatherBaseURL string        `mapstructure:"WEATHER_BASE_URL"`
	RatesAPIKey    string        `mapstructure:"RATES_API_KEY"`
	RatesBaseURL   string        `mapstructure:"RATES_BASE_URL"`
	RequestTimeout time.Duration `mapstructure:"REQUEST_TIMEOUT"`
	CacheTTL       time.Duration `mapstructure:"CACHE_TTL"`
	LogLevel       string        `mapstructure:"LOG_LEVEL"`
}

func Load() (Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetDefault("PORT", "8081")
	v.SetDefault("REQUEST_TIMEOUT", "5s")
	v.SetDefault("CACHE_TTL", "30s")
	v.SetDefault("WEATHER_BASE_URL", "https://api.openweathermap.org/data/2.5")
	v.SetDefault("RATES_BASE_URL", "https://open.er-api.com/v6")
	var c Config
	if err := v.ReadInConfig(); err != nil {
		// ignore if .env not present
	}
	if err := v.Unmarshal(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

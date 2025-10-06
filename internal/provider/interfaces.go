package provider

import (
	"context"
	"github.com/lissoxha/currency/internal/model"
)

type WeatherProvider interface {
	GetWeather(ctx context.Context, city string) (model.WeatherResponse, error)
}

type RatesProvider interface {
	GetRate(ctx context.Context, base string, target string) (model.CurrencyResponse, error)
}

package service

import (
	"context"
	"time"

	"github.com/lissoxha/currency/internal/model"
	"github.com/lissoxha/currency/internal/provider"
)

type Service interface {
	GetWeather(ctx context.Context, city string) (model.WeatherResponse, error)
	GetRate(ctx context.Context, base, target string) (model.CurrencyResponse, error)
	GetReport(ctx context.Context, city, base, target string) (model.ReportResponse, error)
}

type service struct {
	weather provider.WeatherProvider
	rates   provider.RatesProvider
}

func New(weather provider.WeatherProvider, rates provider.RatesProvider) Service {
	return &service{weather: weather, rates: rates}
}

func (s *service) GetWeather(ctx context.Context, city string) (model.WeatherResponse, error) {
	return s.weather.GetWeather(ctx, city)
}

func (s *service) GetRate(ctx context.Context, base, target string) (model.CurrencyResponse, error) {
	return s.rates.GetRate(ctx, base, target)
}

func (s *service) GetReport(ctx context.Context, city, base, target string) (model.ReportResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	weatherCh := make(chan model.WeatherResponse, 1)
	weatherErr := make(chan error, 1)
	ratesCh := make(chan model.CurrencyResponse, 1)
	ratesErr := make(chan error, 1)
	go func() {
		w, err := s.weather.GetWeather(ctx, city)
		if err != nil {
			weatherErr <- err
			return
		}
		weatherCh <- w
	}()
	go func() {
		r, err := s.rates.GetRate(ctx, base, target)
		if err != nil {
			ratesErr <- err
			return
		}
		ratesCh <- r
	}()
	var w model.WeatherResponse
	var r model.CurrencyResponse
	for i := 0; i < 2; i++ {
		select {
		case we := <-weatherErr:
			return model.ReportResponse{}, we
		case re := <-ratesErr:
			return model.ReportResponse{}, re
		case v := <-weatherCh:
			w = v
		case v := <-ratesCh:
			r = v
		case <-ctx.Done():
			return model.ReportResponse{}, ctx.Err()
		}
	}
	return model.ReportResponse{Weather: w, CurrencyRate: r}, nil
}

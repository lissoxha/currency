package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/lissoxha/currency/internal/cache"
	"github.com/lissoxha/currency/internal/config"
	"github.com/lissoxha/currency/internal/model"
)

type OpenWeatherProvider struct {
	cfg   config.Config
	http  *http.Client
	cache *cache.Cache[model.WeatherResponse]
}

func NewOpenWeatherProvider(cfg config.Config, c *cache.Cache[model.WeatherResponse]) *OpenWeatherProvider {
	return &OpenWeatherProvider{cfg: cfg, http: &http.Client{Timeout: cfg.RequestTimeout}, cache: c}
}

func (p *OpenWeatherProvider) GetWeather(ctx context.Context, city string) (model.WeatherResponse, error) {
	if v, ok := p.cache.Get(city); ok {
		return v, nil
	}
	endpoint := fmt.Sprintf("%s/weather", p.cfg.WeatherBaseURL)
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("q", city)
	q.Set("appid", p.cfg.WeatherAPIKey)
	q.Set("units", "metric")
	u.RawQuery = q.Encode()

	var out model.WeatherResponse
	operation := func() error {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		resp, err := p.http.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 500 {
			return fmt.Errorf("upstream %d", resp.StatusCode)
		}
		if resp.StatusCode != 200 {
			return backoff.Permanent(fmt.Errorf("bad status %d", resp.StatusCode))
		}
		var payload struct {
			Name string `json:"name"`
			Main struct {
				Temp float64 `json:"temp"`
			} `json:"main"`
			Weather []struct {
				Main        string `json:"main"`
				Description string `json:"description"`
			} `json:"weather"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return err
		}
		condition := ""
		if len(payload.Weather) > 0 {
			condition = payload.Weather[0].Main
		}
		out = model.WeatherResponse{City: payload.Name, Temperature: payload.Main.Temp, Condition: condition}
		return nil
	}
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 200 * time.Millisecond
	bo.MaxElapsedTime = 2 * time.Second
	if err := backoff.Retry(operation, backoff.WithContext(bo, ctx)); err != nil {
		return model.WeatherResponse{}, err
	}
	p.cache.Set(city, out, p.cfg.CacheTTL)
	return out, nil
}

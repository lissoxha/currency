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

type OpenRatesProvider struct {
	cfg   config.Config
	http  *http.Client
	cache *cache.Cache[model.CurrencyResponse]
}

func NewOpenRatesProvider(cfg config.Config, c *cache.Cache[model.CurrencyResponse]) *OpenRatesProvider {
	return &OpenRatesProvider{cfg: cfg, http: &http.Client{Timeout: cfg.RequestTimeout}, cache: c}
}

func (p *OpenRatesProvider) GetRate(ctx context.Context, base, target string) (model.CurrencyResponse, error) {
	key := base + "_" + target
	if v, ok := p.cache.Get(key); ok {
		return v, nil
	}
	u, _ := url.Parse(fmt.Sprintf("%s/latest/%s", p.cfg.RatesBaseURL, base))
	var out model.CurrencyResponse
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
			Rates    map[string]float64 `json:"rates"`
			BaseCode string             `json:"base_code"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return err
		}
		rate, ok := payload.Rates[target]
		if !ok {
			return backoff.Permanent(fmt.Errorf("rate not found for %s", target))
		}
		out = model.CurrencyResponse{Base: base, Target: target, Rate: rate}
		return nil
	}
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 200 * time.Millisecond
	bo.MaxElapsedTime = 2 * time.Second
	if err := backoff.Retry(operation, backoff.WithContext(bo, ctx)); err != nil {
		return model.CurrencyResponse{}, err
	}
	p.cache.Set(key, out, p.cfg.CacheTTL)
	return out, nil
}

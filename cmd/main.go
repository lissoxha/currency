package main

import (
	"github.com/lissoxha/currency/internal/handler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lissoxha/currency/internal/cache"
	"github.com/lissoxha/currency/internal/config"
	"github.com/lissoxha/currency/internal/model"
	"github.com/lissoxha/currency/internal/provider"
	"github.com/lissoxha/currency/internal/service"
	"go.uber.org/zap"
)

// @title Currency & Weather Aggregator API
// @version 1.0
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// providers
	weatherProv := provider.NewOpenWeatherProvider(cfg, cache.New[model.WeatherResponse]())
	ratesProv := provider.NewOpenRatesProvider(cfg, cache.New[model.CurrencyResponse]())

	svc := service.New(weatherProv, ratesProv)

	r := gin.New()
	r.Use(gin.Recovery())

	h := handler.New(svc)
	h.Register(r)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
	logger.Sugar().Infow("listening", "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Sugar().Fatalw("server error", "err", err)
	}
}

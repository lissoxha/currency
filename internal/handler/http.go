package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lissoxha/currency/internal/service"
)

type Handler struct{ svc service.Service }

func New(svc service.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/weather", h.getWeather)
	r.GET("/currency", h.getCurrency)
	r.GET("/report", h.getReport)
}

func (h *Handler) getWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
		return
	}
	res, err := h.svc.GetWeather(c, city)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getCurrency(c *gin.Context) {
	base := c.Query("base")
	target := c.Query("target")
	if base == "" || target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "base and target are required"})
		return
	}
	res, err := h.svc.GetRate(c, base, target)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getReport(c *gin.Context) {
	city := c.Query("city")
	base := c.Query("base")
	target := c.Query("target")
	if city == "" || base == "" || target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city, base, target required"})
		return
	}
	res, err := h.svc.GetReport(c, city, base, target)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

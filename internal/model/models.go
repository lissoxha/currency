package model

type WeatherResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
}

type CurrencyResponse struct {
	Base   string  `json:"base"`
	Target string  `json:"target"`
	Rate   float64 `json:"rate"`
}

type ReportResponse struct {
	Weather      WeatherResponse  `json:"weather"`
	CurrencyRate CurrencyResponse `json:"currency_rate"`
}

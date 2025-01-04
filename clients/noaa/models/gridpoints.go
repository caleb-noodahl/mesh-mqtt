package models

import "time"

// ForecastResponse represents the response from the forecast endpoint.
type ForecastResponse struct {
	Context []interface{} `json:"@context"`
	Type    string        `json:"type"`
	//Geometry   ForecastGeometry   `json:"geometry"`
	Properties ForecastProperties `json:"properties"`
}

type ForecastGeometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"` //something's wrong with this
}

type ForecastProperties struct {
	Units             string            `json:"units"`
	ForecastGenerator string            `json:"forecastGenerator"`
	GeneratedAt       time.Time         `json:"generatedAt"`
	UpdateTime        time.Time         `json:"updateTime"`
	ValidTimes        string            `json:"validTimes"`
	Elevation         QuantitativeValue `json:"elevation"`
	Periods           []ForecastPeriod  `json:"periods"`
}

type ForecastPeriod struct {
	Number                     int                `json:"number"`
	Name                       string             `json:"name"`
	StartTime                  time.Time          `json:"startTime"`
	EndTime                    time.Time          `json:"endTime"`
	IsDaytime                  bool               `json:"isDaytime"`
	Temperature                int                `json:"temperature"`
	TemperatureUnit            string             `json:"temperatureUnit"`
	TemperatureTrend           string             `json:"temperatureTrend"`
	ProbabilityOfPrecipitation *QuantitativeValue `json:"probabilityOfPrecipitation"`
	WindSpeed                  string             `json:"windSpeed"`
	WindDirection              string             `json:"windDirection"`
	Icon                       string             `json:"icon"`
	ShortForecast              string             `json:"shortForecast"`
	DetailedForecast           string             `json:"detailedForecast"`
}

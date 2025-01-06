package noaa

import (
	"fmt"
	"mesh-mqtt/clients/noaa/models"

	"github.com/go-resty/resty/v2"
)

type NOAAClient struct {
	client *resty.Client
}

func NewNOAAClient(opts ...func(*NOAAClient)) *NOAAClient {
	c := &NOAAClient{
		client: resty.New(),
	}
	c.client.BaseURL = "https://api.weather.gov"
	c.client.SetHeader("Accept", "application/json")
	c.client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "mesh-mqtt,admin@noodahl.com",
	})
	for _, o := range opts {
		o(c)
	}
	return c
}

func (n *NOAAClient) Points(lat, long float64) (*models.PointResponse, error) {
	result := &models.PointResponse{}
	resp, err := n.client.R().
		SetResult(result).
		Get(fmt.Sprintf("%s/points/%v,%v", n.client.BaseURL, lat, long))
	if resp.IsError() {
		return nil, err
	}
	return result, err
}

func (n *NOAAClient) GridPointsForecast(gridID string, gridX, gridY int64) (*models.ForecastResponse, error) {
	result := models.ForecastResponse{}
	resp, err := n.client.R().
		SetResult(&result).
		Get(fmt.Sprintf(
			"%s/gridpoints/%s/%v,%v/forecast",
			n.client.BaseURL,
			gridID,
			gridX,
			gridY),
		)
	if resp.IsError() {
		return nil, err
	}
	return &result, err
}

func (n *NOAAClient) Alerts(lat, long float64) (*models.AlertResponse, error) {
	result := models.AlertResponse{}
	resp, err := n.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("/alerts/active?point=%v,%v", lat, long))
	if resp.IsError() {
		return nil, err
	}
	return &result, err
}

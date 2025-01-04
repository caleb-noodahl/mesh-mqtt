package models

// PointResponse represents the response from the /points/{point} endpoint.
type PointResponse struct {
	Context    []interface{}   `json:"@context"`
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Geometry   Geometry        `json:"geometry"`
	Properties PointProperties `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type PointProperties struct {
	ID                  string           `json:"@id"`
	Type                string           `json:"@type"`
	CWA                 string           `json:"cwa"`
	ForecastOffice      string           `json:"forecastOffice"`
	GridID              string           `json:"gridId"`
	GridX               int              `json:"gridX"`
	GridY               int              `json:"gridY"`
	Forecast            string           `json:"forecast"`
	ForecastHourly      string           `json:"forecastHourly"`
	ForecastGridData    string           `json:"forecastGridData"`
	ObservationStations string           `json:"observationStations"`
	RelativeLocation    RelativeLocation `json:"relativeLocation"`
	ForecastZone        string           `json:"forecastZone"`
	County              string           `json:"county"`
	FireWeatherZone     string           `json:"fireWeatherZone"`
	TimeZone            string           `json:"timeZone"`
	RadarStation        string           `json:"radarStation"`
}

type RelativeLocation struct {
	Type       string             `json:"type"`
	Geometry   Geometry           `json:"geometry"`
	Properties LocationProperties `json:"properties"`
}

type LocationProperties struct {
	City     string            `json:"city"`
	State    string            `json:"state"`
	Distance QuantitativeValue `json:"distance"`
	Bearing  QuantitativeValue `json:"bearing"`
}

type QuantitativeValue struct {
	UnitCode string  `json:"unitCode"`
	Value    float64 `json:"value"`
}

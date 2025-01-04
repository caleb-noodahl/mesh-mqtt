package models

type NodeRecord struct {
	ID        string  `json:"id"`
	Connect   int64   `json:"connect"`
	Next      int64   `json:"next"`
	Last      int64   `json:"last"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	GridID    string  `json:"grid_id"`
	GridX     int     `json:"grid_x"`
	GridY     int     `json:"grid_y"`
}

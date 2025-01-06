package models

// AlertResponse defines the structure for NOAA alert data.
type AlertResponse struct {
	Title   string  `json:"title"`
	Updated string  `json:"updated"`
	Alerts  []Alert `json:"features"`
}

type Alert struct {
	ID         string      `json:"id"`
	Properties AlertDetail `json:"properties"`
}

type AlertDetail struct {
	Headline      string              `json:"headline"`
	Description   string              `json:"description"`
	Instruction   string              `json:"instruction"`
	Event         string              `json:"event"`
	Effective     string              `json:"effective"`
	Expires       string              `json:"expires"`
	Urgency       string              `json:"urgency"`
	Severity      string              `json:"severity"`
	Certainty     string              `json:"certainty"`
	Sender        string              `json:"sender"`
	SenderName    string              `json:"senderName"`
	AreaDesc      string              `json:"areaDesc"`
	Status        string              `json:"status"`
	MessageType   string              `json:"messageType"`
	Category      string              `json:"category"`
	Response      string              `json:"response"`
	Parameters    map[string][]string `json:"parameters"`
	AffectedZones []string            `json:"affectedZones"`
}

package heatmaps

type TrafficFlow struct {
	Latitude         *float64           `json:"latitude,omitempty"`
	Longitude        *float64           `json:"longitude,omitempty"`
	QueryExperiences *[]QueryExperience `json:"queryExperiences,omitempty"`
	SourceIp         *string            `json:"sourceIp,omitempty"`
}

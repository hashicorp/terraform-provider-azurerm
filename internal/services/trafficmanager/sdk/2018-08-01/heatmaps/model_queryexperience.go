package heatmaps

type QueryExperience struct {
	EndpointId int64    `json:"endpointId"`
	Latency    *float64 `json:"latency,omitempty"`
	QueryCount int64    `json:"queryCount"`
}

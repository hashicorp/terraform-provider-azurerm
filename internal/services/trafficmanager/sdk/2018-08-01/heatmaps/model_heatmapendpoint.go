package heatmaps

type HeatMapEndpoint struct {
	EndpointId *int64  `json:"endpointId,omitempty"`
	ResourceId *string `json:"resourceId,omitempty"`
}

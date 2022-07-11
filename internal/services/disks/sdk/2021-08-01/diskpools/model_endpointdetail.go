package diskpools

type EndpointDetail struct {
	IpAddress    *string  `json:"ipAddress,omitempty"`
	IsAccessible *bool    `json:"isAccessible,omitempty"`
	Latency      *float64 `json:"latency,omitempty"`
	Port         *int64   `json:"port,omitempty"`
}

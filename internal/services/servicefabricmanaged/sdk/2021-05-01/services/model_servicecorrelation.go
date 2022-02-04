package services

type ServiceCorrelation struct {
	Scheme      ServiceCorrelationScheme `json:"scheme"`
	ServiceName string                   `json:"serviceName"`
}

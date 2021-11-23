package services

type ServiceLoadMetric struct {
	DefaultLoad          *int64                   `json:"defaultLoad,omitempty"`
	Name                 string                   `json:"name"`
	PrimaryDefaultLoad   *int64                   `json:"primaryDefaultLoad,omitempty"`
	SecondaryDefaultLoad *int64                   `json:"secondaryDefaultLoad,omitempty"`
	Weight               *ServiceLoadMetricWeight `json:"weight,omitempty"`
}

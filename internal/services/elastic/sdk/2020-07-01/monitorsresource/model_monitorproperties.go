package monitorsresource

type MonitorProperties struct {
	ElasticProperties       *ElasticProperties       `json:"elasticProperties,omitempty"`
	LiftrResourceCategory   *LiftrResourceCategories `json:"liftrResourceCategory,omitempty"`
	LiftrResourcePreference *int64                   `json:"liftrResourcePreference,omitempty"`
	MonitoringStatus        *MonitoringStatus        `json:"monitoringStatus,omitempty"`
	ProvisioningState       *ProvisioningState       `json:"provisioningState,omitempty"`
	UserInfo                *UserInfo                `json:"userInfo,omitempty"`
}

package cognitiveservicesaccounts

type Usage struct {
	CurrentValue  *float64          `json:"currentValue,omitempty"`
	Limit         *float64          `json:"limit,omitempty"`
	Name          *MetricName       `json:"name,omitempty"`
	NextResetTime *string           `json:"nextResetTime,omitempty"`
	QuotaPeriod   *string           `json:"quotaPeriod,omitempty"`
	Status        *QuotaUsageStatus `json:"status,omitempty"`
	Unit          *UnitType         `json:"unit,omitempty"`
}

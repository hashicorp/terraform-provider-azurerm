package alerts

type AlertPropertiesDetails struct {
	Amount              *float64            `json:"amount,omitempty"`
	ContactEmails       *[]string           `json:"contactEmails,omitempty"`
	ContactGroups       *[]string           `json:"contactGroups,omitempty"`
	ContactRoles        *[]string           `json:"contactRoles,omitempty"`
	CurrentSpend        *float64            `json:"currentSpend,omitempty"`
	MeterFilter         *[]interface{}      `json:"meterFilter,omitempty"`
	Operator            *AlertOperator      `json:"operator,omitempty"`
	OverridingAlert     *string             `json:"overridingAlert,omitempty"`
	PeriodStartDate     *string             `json:"periodStartDate,omitempty"`
	ResourceFilter      *[]interface{}      `json:"resourceFilter,omitempty"`
	ResourceGroupFilter *[]interface{}      `json:"resourceGroupFilter,omitempty"`
	TagFilter           *interface{}        `json:"tagFilter,omitempty"`
	Threshold           *float64            `json:"threshold,omitempty"`
	TimeGrainType       *AlertTimeGrainType `json:"timeGrainType,omitempty"`
	TriggeredBy         *string             `json:"triggeredBy,omitempty"`
	Unit                *string             `json:"unit,omitempty"`
}

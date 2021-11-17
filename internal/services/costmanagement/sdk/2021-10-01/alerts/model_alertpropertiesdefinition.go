package alerts

type AlertPropertiesDefinition struct {
	Category *AlertCategory `json:"category,omitempty"`
	Criteria *AlertCriteria `json:"criteria,omitempty"`
	Type     *AlertType     `json:"type,omitempty"`
}

package realusermetrics

type UserMetricsModel struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties *UserMetricsProperties `json:"properties,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

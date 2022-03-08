package heatmaps

type HeatMapModel struct {
	Id         *string            `json:"id,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties *HeatMapProperties `json:"properties,omitempty"`
	Type       *string            `json:"type,omitempty"`
}

package views

type KpiProperties struct {
	Enabled *bool        `json:"enabled,omitempty"`
	Id      *string      `json:"id,omitempty"`
	Type    *KpiTypeType `json:"type,omitempty"`
}

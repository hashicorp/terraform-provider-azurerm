package views

type PivotProperties struct {
	Name *string        `json:"name,omitempty"`
	Type *PivotTypeType `json:"type,omitempty"`
}

package topics

type EventTypeProperties struct {
	Description    *string `json:"description,omitempty"`
	DisplayName    *string `json:"displayName,omitempty"`
	IsInDefaultSet *bool   `json:"isInDefaultSet,omitempty"`
	SchemaUrl      *string `json:"schemaUrl,omitempty"`
}

package topics

type JsonFieldWithDefault struct {
	DefaultValue *string `json:"defaultValue,omitempty"`
	SourceField  *string `json:"sourceField,omitempty"`
}

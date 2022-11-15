package frontdoors

type RulesEngine struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties *RulesEngineProperties `json:"properties,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

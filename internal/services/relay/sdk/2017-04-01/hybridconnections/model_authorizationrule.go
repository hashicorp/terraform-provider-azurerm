package hybridconnections

type AuthorizationRule struct {
	Id         *string                     `json:"id,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	Properties AuthorizationRuleProperties `json:"properties"`
	Type       *string                     `json:"type,omitempty"`
}

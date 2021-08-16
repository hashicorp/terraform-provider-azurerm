package authorizationruleseventhubs

type AuthorizationRule struct {
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties *AuthorizationRuleProperties `json:"properties,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}

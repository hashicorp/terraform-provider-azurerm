package authorizationruleseventhubs

type AuthorizationRule struct {
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties *AuthorizationRuleProperties `json:"properties,omitempty"`
	SystemData *SystemData                  `json:"systemData,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}

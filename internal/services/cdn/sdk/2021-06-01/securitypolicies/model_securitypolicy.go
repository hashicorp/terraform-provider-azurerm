package securitypolicies

type SecurityPolicy struct {
	Id         *string                   `json:"id,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *SecurityPolicyProperties `json:"properties,omitempty"`
	SystemData *SystemData               `json:"systemData,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}

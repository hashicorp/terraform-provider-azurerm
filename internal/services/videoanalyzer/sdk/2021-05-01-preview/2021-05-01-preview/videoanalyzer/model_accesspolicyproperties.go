package videoanalyzer

type AccessPolicyProperties struct {
	Authentication *AuthenticationBase `json:"authentication,omitempty"`
	Role           *AccessPolicyRole   `json:"role,omitempty"`
}

package webapplicationfirewallpolicies

type PolicySettings struct {
	CustomBlockResponseBody       *string             `json:"customBlockResponseBody,omitempty"`
	CustomBlockResponseStatusCode *int64              `json:"customBlockResponseStatusCode,omitempty"`
	EnabledState                  *PolicyEnabledState `json:"enabledState,omitempty"`
	Mode                          *PolicyMode         `json:"mode,omitempty"`
	RedirectUrl                   *string             `json:"redirectUrl,omitempty"`
}

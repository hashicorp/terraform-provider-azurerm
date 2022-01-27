package webapplicationfirewallpolicies

type PolicySettings struct {
	DefaultCustomBlockResponseBody       *string                               `json:"defaultCustomBlockResponseBody,omitempty"`
	DefaultCustomBlockResponseStatusCode *DefaultCustomBlockResponseStatusCode `json:"defaultCustomBlockResponseStatusCode,omitempty"`
	DefaultRedirectUrl                   *string                               `json:"defaultRedirectUrl,omitempty"`
	EnabledState                         *PolicyEnabledState                   `json:"enabledState,omitempty"`
	Mode                                 *PolicyMode                           `json:"mode,omitempty"`
}

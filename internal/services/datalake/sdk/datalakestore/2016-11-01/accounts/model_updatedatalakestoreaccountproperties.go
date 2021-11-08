package accounts

type UpdateDataLakeStoreAccountProperties struct {
	DefaultGroup           *string                                          `json:"defaultGroup,omitempty"`
	EncryptionConfig       *UpdateEncryptionConfig                          `json:"encryptionConfig,omitempty"`
	FirewallAllowAzureIps  *FirewallAllowAzureIpsState                      `json:"firewallAllowAzureIps,omitempty"`
	FirewallRules          *[]UpdateFirewallRuleWithAccountParameters       `json:"firewallRules,omitempty"`
	FirewallState          *FirewallState                                   `json:"firewallState,omitempty"`
	NewTier                *TierType                                        `json:"newTier,omitempty"`
	TrustedIdProviderState *TrustedIdProviderState                          `json:"trustedIdProviderState,omitempty"`
	TrustedIdProviders     *[]UpdateTrustedIdProviderWithAccountParameters  `json:"trustedIdProviders,omitempty"`
	VirtualNetworkRules    *[]UpdateVirtualNetworkRuleWithAccountParameters `json:"virtualNetworkRules,omitempty"`
}

package accounts

type CreateDataLakeStoreAccountProperties struct {
	DefaultGroup           *string                                          `json:"defaultGroup,omitempty"`
	EncryptionConfig       *EncryptionConfig                                `json:"encryptionConfig,omitempty"`
	EncryptionState        *EncryptionState                                 `json:"encryptionState,omitempty"`
	FirewallAllowAzureIps  *FirewallAllowAzureIpsState                      `json:"firewallAllowAzureIps,omitempty"`
	FirewallRules          *[]CreateFirewallRuleWithAccountParameters       `json:"firewallRules,omitempty"`
	FirewallState          *FirewallState                                   `json:"firewallState,omitempty"`
	NewTier                *TierType                                        `json:"newTier,omitempty"`
	TrustedIdProviderState *TrustedIdProviderState                          `json:"trustedIdProviderState,omitempty"`
	TrustedIdProviders     *[]CreateTrustedIdProviderWithAccountParameters  `json:"trustedIdProviders,omitempty"`
	VirtualNetworkRules    *[]CreateVirtualNetworkRuleWithAccountParameters `json:"virtualNetworkRules,omitempty"`
}

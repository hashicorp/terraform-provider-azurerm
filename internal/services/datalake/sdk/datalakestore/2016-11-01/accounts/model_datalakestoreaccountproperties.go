package accounts

type DataLakeStoreAccountProperties struct {
	AccountId                   *string                      `json:"accountId,omitempty"`
	CreationTime                *string                      `json:"creationTime,omitempty"`
	CurrentTier                 *TierType                    `json:"currentTier,omitempty"`
	DefaultGroup                *string                      `json:"defaultGroup,omitempty"`
	EncryptionConfig            *EncryptionConfig            `json:"encryptionConfig,omitempty"`
	EncryptionProvisioningState *EncryptionProvisioningState `json:"encryptionProvisioningState,omitempty"`
	EncryptionState             *EncryptionState             `json:"encryptionState,omitempty"`
	Endpoint                    *string                      `json:"endpoint,omitempty"`
	FirewallAllowAzureIps       *FirewallAllowAzureIpsState  `json:"firewallAllowAzureIps,omitempty"`
	FirewallRules               *[]FirewallRule              `json:"firewallRules,omitempty"`
	FirewallState               *FirewallState               `json:"firewallState,omitempty"`
	LastModifiedTime            *string                      `json:"lastModifiedTime,omitempty"`
	NewTier                     *TierType                    `json:"newTier,omitempty"`
	ProvisioningState           *DataLakeStoreAccountStatus  `json:"provisioningState,omitempty"`
	State                       *DataLakeStoreAccountState   `json:"state,omitempty"`
	TrustedIdProviderState      *TrustedIdProviderState      `json:"trustedIdProviderState,omitempty"`
	TrustedIdProviders          *[]TrustedIdProvider         `json:"trustedIdProviders,omitempty"`
	VirtualNetworkRules         *[]VirtualNetworkRule        `json:"virtualNetworkRules,omitempty"`
}

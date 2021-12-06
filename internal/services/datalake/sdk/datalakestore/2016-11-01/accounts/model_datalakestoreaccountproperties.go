package accounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

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

func (o DataLakeStoreAccountProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o DataLakeStoreAccountProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o DataLakeStoreAccountProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o DataLakeStoreAccountProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

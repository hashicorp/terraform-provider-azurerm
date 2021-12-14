package accounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type DataLakeAnalyticsAccountProperties struct {
	AccountId                    *string                            `json:"accountId,omitempty"`
	ComputePolicies              *[]ComputePolicy                   `json:"computePolicies,omitempty"`
	CreationTime                 *string                            `json:"creationTime,omitempty"`
	CurrentTier                  *TierType                          `json:"currentTier,omitempty"`
	DataLakeStoreAccounts        *[]DataLakeStoreAccountInformation `json:"dataLakeStoreAccounts,omitempty"`
	DebugDataAccessLevel         *DebugDataAccessLevel              `json:"debugDataAccessLevel,omitempty"`
	DefaultDataLakeStoreAccount  *string                            `json:"defaultDataLakeStoreAccount,omitempty"`
	Endpoint                     *string                            `json:"endpoint,omitempty"`
	FirewallAllowAzureIps        *FirewallAllowAzureIpsState        `json:"firewallAllowAzureIps,omitempty"`
	FirewallRules                *[]FirewallRule                    `json:"firewallRules,omitempty"`
	FirewallState                *FirewallState                     `json:"firewallState,omitempty"`
	HiveMetastores               *[]HiveMetastore                   `json:"hiveMetastores,omitempty"`
	LastModifiedTime             *string                            `json:"lastModifiedTime,omitempty"`
	MaxActiveJobCountPerUser     *int64                             `json:"maxActiveJobCountPerUser,omitempty"`
	MaxDegreeOfParallelism       *int64                             `json:"maxDegreeOfParallelism,omitempty"`
	MaxDegreeOfParallelismPerJob *int64                             `json:"maxDegreeOfParallelismPerJob,omitempty"`
	MaxJobCount                  *int64                             `json:"maxJobCount,omitempty"`
	MaxJobRunningTimeInMin       *int64                             `json:"maxJobRunningTimeInMin,omitempty"`
	MaxQueuedJobCountPerUser     *int64                             `json:"maxQueuedJobCountPerUser,omitempty"`
	MinPriorityPerJob            *int64                             `json:"minPriorityPerJob,omitempty"`
	NewTier                      *TierType                          `json:"newTier,omitempty"`
	ProvisioningState            *DataLakeAnalyticsAccountStatus    `json:"provisioningState,omitempty"`
	PublicDataLakeStoreAccounts  *[]DataLakeStoreAccountInformation `json:"publicDataLakeStoreAccounts,omitempty"`
	QueryStoreRetention          *int64                             `json:"queryStoreRetention,omitempty"`
	State                        *DataLakeAnalyticsAccountState     `json:"state,omitempty"`
	StorageAccounts              *[]StorageAccountInformation       `json:"storageAccounts,omitempty"`
	SystemMaxDegreeOfParallelism *int64                             `json:"systemMaxDegreeOfParallelism,omitempty"`
	SystemMaxJobCount            *int64                             `json:"systemMaxJobCount,omitempty"`
	VirtualNetworkRules          *[]VirtualNetworkRule              `json:"virtualNetworkRules,omitempty"`
}

func (o DataLakeAnalyticsAccountProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o DataLakeAnalyticsAccountProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o DataLakeAnalyticsAccountProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o DataLakeAnalyticsAccountProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

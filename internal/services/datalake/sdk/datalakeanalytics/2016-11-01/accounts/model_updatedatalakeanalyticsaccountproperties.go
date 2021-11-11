package accounts

type UpdateDataLakeAnalyticsAccountProperties struct {
	ComputePolicies              *[]UpdateComputePolicyWithAccountParameters  `json:"computePolicies,omitempty"`
	DataLakeStoreAccounts        *[]UpdateDataLakeStoreWithAccountParameters  `json:"dataLakeStoreAccounts,omitempty"`
	FirewallAllowAzureIps        *FirewallAllowAzureIpsState                  `json:"firewallAllowAzureIps,omitempty"`
	FirewallRules                *[]UpdateFirewallRuleWithAccountParameters   `json:"firewallRules,omitempty"`
	FirewallState                *FirewallState                               `json:"firewallState,omitempty"`
	MaxDegreeOfParallelism       *int64                                       `json:"maxDegreeOfParallelism,omitempty"`
	MaxDegreeOfParallelismPerJob *int64                                       `json:"maxDegreeOfParallelismPerJob,omitempty"`
	MaxJobCount                  *int64                                       `json:"maxJobCount,omitempty"`
	MinPriorityPerJob            *int64                                       `json:"minPriorityPerJob,omitempty"`
	NewTier                      *TierType                                    `json:"newTier,omitempty"`
	QueryStoreRetention          *int64                                       `json:"queryStoreRetention,omitempty"`
	StorageAccounts              *[]UpdateStorageAccountWithAccountParameters `json:"storageAccounts,omitempty"`
}

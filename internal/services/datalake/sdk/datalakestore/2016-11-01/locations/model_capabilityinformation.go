package locations

type CapabilityInformation struct {
	AccountCount    *int64             `json:"accountCount,omitempty"`
	MaxAccountCount *int64             `json:"maxAccountCount,omitempty"`
	MigrationState  *bool              `json:"migrationState,omitempty"`
	State           *SubscriptionState `json:"state,omitempty"`
	SubscriptionId  *string            `json:"subscriptionId,omitempty"`
}

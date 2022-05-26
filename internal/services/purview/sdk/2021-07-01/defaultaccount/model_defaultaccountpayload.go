package defaultaccount

type DefaultAccountPayload struct {
	AccountName       *string    `json:"accountName,omitempty"`
	ResourceGroupName *string    `json:"resourceGroupName,omitempty"`
	Scope             *string    `json:"scope,omitempty"`
	ScopeTenantId     *string    `json:"scopeTenantId,omitempty"`
	ScopeType         *ScopeType `json:"scopeType,omitempty"`
	SubscriptionId    *string    `json:"subscriptionId,omitempty"`
}

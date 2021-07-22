package disasterrecoveryconfigs

type ArmDisasterRecoveryProperties struct {
	AlternateName                     *string               `json:"alternateName,omitempty"`
	PartnerNamespace                  *string               `json:"partnerNamespace,omitempty"`
	PendingReplicationOperationsCount *int64                `json:"pendingReplicationOperationsCount,omitempty"`
	ProvisioningState                 *ProvisioningStateDR  `json:"provisioningState,omitempty"`
	Role                              *RoleDisasterRecovery `json:"role,omitempty"`
}

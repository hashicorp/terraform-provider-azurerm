package disasterrecoveryconfigs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArmDisasterRecoveryProperties struct {
	AlternateName                     *string               `json:"alternateName,omitempty"`
	PartnerNamespace                  *string               `json:"partnerNamespace,omitempty"`
	PendingReplicationOperationsCount *int64                `json:"pendingReplicationOperationsCount,omitempty"`
	ProvisioningState                 *ProvisioningStateDR  `json:"provisioningState,omitempty"`
	Role                              *RoleDisasterRecovery `json:"role,omitempty"`
}

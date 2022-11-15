package resourceguards

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuard struct {
	AllowAutoApprovals                  *bool                     `json:"allowAutoApprovals,omitempty"`
	Description                         *string                   `json:"description,omitempty"`
	ProvisioningState                   *ProvisioningState        `json:"provisioningState,omitempty"`
	ResourceGuardOperations             *[]ResourceGuardOperation `json:"resourceGuardOperations,omitempty"`
	VaultCriticalOperationExclusionList *[]string                 `json:"vaultCriticalOperationExclusionList,omitempty"`
}

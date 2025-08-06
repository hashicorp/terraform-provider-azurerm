package updateruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateRunProperties struct {
	ManagedClusterUpdate ManagedClusterUpdate        `json:"managedClusterUpdate"`
	ProvisioningState    *UpdateRunProvisioningState `json:"provisioningState,omitempty"`
	Status               *UpdateRunStatus            `json:"status,omitempty"`
	Strategy             *UpdateRunStrategy          `json:"strategy,omitempty"`
	UpdateStrategyId     *string                     `json:"updateStrategyId,omitempty"`
}

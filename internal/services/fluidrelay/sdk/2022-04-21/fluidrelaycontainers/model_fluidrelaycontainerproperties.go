package fluidrelaycontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluidRelayContainerProperties struct {
	FrsContainerId    *string            `json:"frsContainerId,omitempty"`
	FrsTenantId       *string            `json:"frsTenantId,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

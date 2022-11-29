package fluidrelayservers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluidRelayServerProperties struct {
	Encryption          *EncryptionProperties `json:"encryption,omitempty"`
	FluidRelayEndpoints *FluidRelayEndpoints  `json:"fluidRelayEndpoints,omitempty"`
	FrsTenantId         *string               `json:"frsTenantId,omitempty"`
	ProvisioningState   *ProvisioningState    `json:"provisioningState,omitempty"`
	Storagesku          *StorageSKU           `json:"storagesku,omitempty"`
}

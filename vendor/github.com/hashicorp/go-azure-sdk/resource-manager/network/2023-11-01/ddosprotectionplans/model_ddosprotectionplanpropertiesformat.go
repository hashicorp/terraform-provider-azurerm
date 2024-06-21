package ddosprotectionplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosProtectionPlanPropertiesFormat struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	PublicIPAddresses *[]SubResource     `json:"publicIPAddresses,omitempty"`
	ResourceGuid      *string            `json:"resourceGuid,omitempty"`
	VirtualNetworks   *[]SubResource     `json:"virtualNetworks,omitempty"`
}

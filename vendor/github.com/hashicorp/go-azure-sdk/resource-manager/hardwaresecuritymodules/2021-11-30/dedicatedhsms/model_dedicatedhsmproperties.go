package dedicatedhsms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHsmProperties struct {
	ManagementNetworkProfile *NetworkProfile `json:"managementNetworkProfile,omitempty"`
	NetworkProfile           *NetworkProfile `json:"networkProfile,omitempty"`
	ProvisioningState        *JsonWebKeyType `json:"provisioningState,omitempty"`
	StampId                  *string         `json:"stampId,omitempty"`
	StatusMessage            *string         `json:"statusMessage,omitempty"`
}

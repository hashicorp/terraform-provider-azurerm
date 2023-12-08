package devcenters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevCenterProperties struct {
	DevCenterUri      *string            `json:"devCenterUri,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

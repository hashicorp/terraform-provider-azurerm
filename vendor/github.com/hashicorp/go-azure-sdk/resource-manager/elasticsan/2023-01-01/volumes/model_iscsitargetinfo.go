package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IscsiTargetInfo struct {
	ProvisioningState    *ProvisioningStates `json:"provisioningState,omitempty"`
	Status               *OperationalStatus  `json:"status,omitempty"`
	TargetIqn            *string             `json:"targetIqn,omitempty"`
	TargetPortalHostname *string             `json:"targetPortalHostname,omitempty"`
	TargetPortalPort     *int64              `json:"targetPortalPort,omitempty"`
}

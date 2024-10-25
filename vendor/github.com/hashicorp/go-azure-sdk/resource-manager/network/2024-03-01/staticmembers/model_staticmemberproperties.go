package staticmembers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticMemberProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Region            *string            `json:"region,omitempty"`
	ResourceId        *string            `json:"resourceId,omitempty"`
}

package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPatchProperties struct {
	AadClientId             *string                   `json:"aadClientId,omitempty"`
	AadTenantId             *string                   `json:"aadTenantId,omitempty"`
	CloudManagementEndpoint *string                   `json:"cloudManagementEndpoint,omitempty"`
	DesiredProperties       *ClusterDesiredProperties `json:"desiredProperties,omitempty"`
}

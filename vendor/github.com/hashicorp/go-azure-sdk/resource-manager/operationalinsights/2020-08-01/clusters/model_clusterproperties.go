package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	ClusterId          *string              `json:"clusterId,omitempty"`
	KeyVaultProperties *KeyVaultProperties  `json:"keyVaultProperties,omitempty"`
	NextLink           *string              `json:"nextLink,omitempty"`
	ProvisioningState  *ClusterEntityStatus `json:"provisioningState,omitempty"`
}

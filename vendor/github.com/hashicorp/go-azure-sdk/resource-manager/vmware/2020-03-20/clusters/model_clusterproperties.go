package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	ClusterId         *int64                    `json:"clusterId,omitempty"`
	ClusterSize       int64                     `json:"clusterSize"`
	Hosts             *[]string                 `json:"hosts,omitempty"`
	ProvisioningState *ClusterProvisioningState `json:"provisioningState,omitempty"`
}

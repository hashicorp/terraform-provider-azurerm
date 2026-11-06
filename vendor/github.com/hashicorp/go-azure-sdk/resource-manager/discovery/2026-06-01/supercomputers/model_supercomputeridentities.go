package supercomputers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupercomputerIdentities struct {
	ClusterIdentity    Identity                         `json:"clusterIdentity"`
	KubeletIdentity    Identity                         `json:"kubeletIdentity"`
	WorkloadIdentities *map[string]UserAssignedIdentity `json:"workloadIdentities,omitempty"`
}

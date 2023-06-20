package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraClusterPublicStatusDataCentersInlined struct {
	Name      *string                                                       `json:"name,omitempty"`
	Nodes     *[]CassandraClusterPublicStatusDataCentersInlinedNodesInlined `json:"nodes,omitempty"`
	SeedNodes *[]string                                                     `json:"seedNodes,omitempty"`
}

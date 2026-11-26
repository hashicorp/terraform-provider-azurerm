package openapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraClusterPublicStatusDataCentersItem struct {
	Name      *string                                                                                              `json:"name,omitempty"`
	Nodes     *[]ComponentsM9L909SchemasCassandraclusterpublicstatusPropertiesDatacentersItemsPropertiesNodesItems `json:"nodes,omitempty"`
	SeedNodes *[]string                                                                                            `json:"seedNodes,omitempty"`
}

package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SQLServer struct {
	ClusterName *string `json:"clusterName,omitempty"`
	Clustered   *string `json:"clustered,omitempty"`
	Edition     *string `json:"edition,omitempty"`
	Name        *string `json:"name,omitempty"`
	ServicePack *string `json:"servicePack,omitempty"`
	Version     *string `json:"version,omitempty"`
}

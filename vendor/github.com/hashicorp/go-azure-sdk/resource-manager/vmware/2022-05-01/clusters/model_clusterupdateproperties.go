package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterUpdateProperties struct {
	ClusterSize *int64    `json:"clusterSize,omitempty"`
	Hosts       *[]string `json:"hosts,omitempty"`
}

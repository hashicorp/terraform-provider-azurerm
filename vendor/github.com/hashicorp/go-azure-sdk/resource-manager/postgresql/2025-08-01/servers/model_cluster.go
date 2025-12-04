package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Cluster struct {
	ClusterSize         *int64  `json:"clusterSize,omitempty"`
	DefaultDatabaseName *string `json:"defaultDatabaseName,omitempty"`
}

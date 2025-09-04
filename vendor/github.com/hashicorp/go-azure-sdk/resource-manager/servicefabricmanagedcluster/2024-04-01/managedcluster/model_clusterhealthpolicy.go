package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterHealthPolicy struct {
	MaxPercentUnhealthyApplications int64 `json:"maxPercentUnhealthyApplications"`
	MaxPercentUnhealthyNodes        int64 `json:"maxPercentUnhealthyNodes"`
}

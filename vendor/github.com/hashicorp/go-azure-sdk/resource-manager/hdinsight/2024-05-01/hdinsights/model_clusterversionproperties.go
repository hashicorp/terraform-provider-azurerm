package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterVersionProperties struct {
	ClusterPoolVersion *string                               `json:"clusterPoolVersion,omitempty"`
	ClusterType        *string                               `json:"clusterType,omitempty"`
	ClusterVersion     *string                               `json:"clusterVersion,omitempty"`
	Components         *[]ClusterComponentsComponentsInlined `json:"components,omitempty"`
	IsPreview          *bool                                 `json:"isPreview,omitempty"`
	OssVersion         *string                               `json:"ossVersion,omitempty"`
}

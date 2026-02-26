package scclusterrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	Metadata *SCMetadataEntity    `json:"metadata,omitempty"`
	Spec     *SCClusterSpecEntity `json:"spec,omitempty"`
	Status   *ClusterStatusEntity `json:"status,omitempty"`
}

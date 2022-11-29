package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPatch struct {
	Properties *ClusterPatchProperties `json:"properties,omitempty"`
	Sku        *ClusterSku             `json:"sku,omitempty"`
	Tags       *map[string]string      `json:"tags,omitempty"`
}

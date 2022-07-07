package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Cluster struct {
	Id         *string           `json:"id,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Properties ClusterProperties `json:"properties"`
	Sku        Sku               `json:"sku"`
	Type       *string           `json:"type,omitempty"`
}

package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedCluster struct {
	Etag       *string                   `json:"etag,omitempty"`
	Id         *string                   `json:"id,omitempty"`
	Location   string                    `json:"location"`
	Name       *string                   `json:"name,omitempty"`
	Properties *ManagedClusterProperties `json:"properties,omitempty"`
	Sku        *Sku                      `json:"sku,omitempty"`
	SystemData *SystemData               `json:"systemData,omitempty"`
	Tags       *map[string]string        `json:"tags,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}

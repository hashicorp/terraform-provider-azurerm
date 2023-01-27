package diskpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskPoolUpdate struct {
	ManagedBy         *string                  `json:"managedBy,omitempty"`
	ManagedByExtended *[]string                `json:"managedByExtended,omitempty"`
	Properties        DiskPoolUpdateProperties `json:"properties"`
	Sku               *Sku                     `json:"sku,omitempty"`
	Tags              *map[string]string       `json:"tags,omitempty"`
}

package snapshots

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Snapshot struct {
	ExtendedLocation *edgezones.Model    `json:"extendedLocation,omitempty"`
	Id               *string             `json:"id,omitempty"`
	Location         string              `json:"location"`
	ManagedBy        *string             `json:"managedBy,omitempty"`
	Name             *string             `json:"name,omitempty"`
	Properties       *SnapshotProperties `json:"properties,omitempty"`
	Sku              *SnapshotSku        `json:"sku,omitempty"`
	Tags             *map[string]string  `json:"tags,omitempty"`
	Type             *string             `json:"type,omitempty"`
}

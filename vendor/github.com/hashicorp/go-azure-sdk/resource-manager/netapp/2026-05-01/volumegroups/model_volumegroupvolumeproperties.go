package volumegroups

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeGroupVolumeProperties struct {
	Id         *string            `json:"id,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties VolumeProperties   `json:"properties"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
	Zones      *zones.Schema      `json:"zones,omitempty"`
}

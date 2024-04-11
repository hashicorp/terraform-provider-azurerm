package dedicatedhostgroups

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostGroupUpdate struct {
	Properties *DedicatedHostGroupProperties `json:"properties,omitempty"`
	Tags       *map[string]string            `json:"tags,omitempty"`
	Zones      *zones.Schema                 `json:"zones,omitempty"`
}

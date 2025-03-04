package virtualmachinescalesets

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSet struct {
	Etag             *string                            `json:"etag,omitempty"`
	ExtendedLocation *edgezones.Model                   `json:"extendedLocation,omitempty"`
	Id               *string                            `json:"id,omitempty"`
	Identity         *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Location         string                             `json:"location"`
	Name             *string                            `json:"name,omitempty"`
	Plan             *Plan                              `json:"plan,omitempty"`
	Properties       *VirtualMachineScaleSetProperties  `json:"properties,omitempty"`
	Sku              *Sku                               `json:"sku,omitempty"`
	Tags             *map[string]string                 `json:"tags,omitempty"`
	Type             *string                            `json:"type,omitempty"`
	Zones            *zones.Schema                      `json:"zones,omitempty"`
}

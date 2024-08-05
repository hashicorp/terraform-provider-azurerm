package virtualmachinescalesetvms

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetVM struct {
	Etag       *string                             `json:"etag,omitempty"`
	Id         *string                             `json:"id,omitempty"`
	Identity   *identity.SystemAndUserAssignedMap  `json:"identity,omitempty"`
	InstanceId *string                             `json:"instanceId,omitempty"`
	Location   string                              `json:"location"`
	Name       *string                             `json:"name,omitempty"`
	Plan       *Plan                               `json:"plan,omitempty"`
	Properties *VirtualMachineScaleSetVMProperties `json:"properties,omitempty"`
	Resources  *[]VirtualMachineExtension          `json:"resources,omitempty"`
	Sku        *Sku                                `json:"sku,omitempty"`
	Tags       *map[string]string                  `json:"tags,omitempty"`
	Type       *string                             `json:"type,omitempty"`
	Zones      *zones.Schema                       `json:"zones,omitempty"`
}

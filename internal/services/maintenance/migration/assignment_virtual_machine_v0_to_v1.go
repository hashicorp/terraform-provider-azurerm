// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = AssignmentVirtualMachineV0ToV1{}

type AssignmentVirtualMachineV0ToV1 struct {
}

func (AssignmentVirtualMachineV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"maintenance_configuration_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"virtual_machine_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (AssignmentVirtualMachineV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := configurationassignments.ParseScopedConfigurationAssignmentIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, fmt.Errorf("parsing old id %q: %+v", oldIdRaw, err)
		}

		virtualMachineId, err := virtualmachines.ParseVirtualMachineIDInsensitively(oldId.Scope)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a virtual machine id: %+v", oldId.Scope, err)
		}

		newId := configurationassignments.NewScopedConfigurationAssignmentID(virtualMachineId.ID(), oldId.ConfigurationAssignmentName)
		newIdRaw := newId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newIdRaw)
		rawState["id"] = newIdRaw

		return rawState, nil
	}
}

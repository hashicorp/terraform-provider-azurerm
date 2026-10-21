// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_virtual_machine

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/resources/maintenance_assignment_virtual_machine/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmMaintenanceAssignmentVirtualMachineCreate,
		Read:   resourceArmMaintenanceAssignmentVirtualMachineRead,
		Delete: resourceArmMaintenanceAssignmentVirtualMachineDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			parsed, err := configurationassignments.ParseScopedConfigurationAssignmentID(id)
			if err != nil {
				return err
			}
			if _, err := virtualmachines.ParseVirtualMachineID(parsed.Scope); err != nil {
				return fmt.Errorf("parsing %q as Virtual Machine ID: %+v", parsed.Scope, err)
			}
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AssignmentVirtualMachineV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: maintenanceAssignmentVirtualMachineSchema(),
	}
}

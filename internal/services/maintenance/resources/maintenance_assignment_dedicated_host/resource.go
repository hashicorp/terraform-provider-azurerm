// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_dedicated_host

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/resources/maintenance_assignment_dedicated_host/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmMaintenanceAssignmentDedicatedHostCreate,
		Read:   resourceArmMaintenanceAssignmentDedicatedHostRead,
		Delete: resourceArmMaintenanceAssignmentDedicatedHostDelete,

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

			if _, err := commonids.ParseDedicatedHostID(parsed.Scope); err != nil {
				return fmt.Errorf("parsing %q as a Dedicated Host ID: %+v", parsed.Scope, err)
			}

			return nil
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AssignmentDedicatedHostV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: maintenanceAssignmentDedicatedHostSchema(),
	}
}

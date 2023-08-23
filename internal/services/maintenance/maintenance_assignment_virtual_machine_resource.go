// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmMaintenanceAssignmentVirtualMachine() *pluginsdk.Resource {
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

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.Location(),

			"maintenance_configuration_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     maintenanceconfigurations.ValidateMaintenanceConfigurationID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualmachines.ValidateVirtualMachineID,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentVirtualMachineCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineId, err := virtualmachines.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}

	configurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Get("maintenance_configuration_id").(string))
	if err != nil {
		return err
	}

	id := configurationassignments.NewScopedConfigurationAssignmentID(virtualMachineId.ID(), configurationId.MaintenanceConfigurationName)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine", id.ID())
	}

	// set assignment name to configuration name
	assignmentName := configurationId.MaintenanceConfigurationName
	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Name:     utils.String(assignmentName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: utils.String(configurationId.ID()),
			ResourceId:                 utils.String(virtualMachineId.ID()),
		},
	}

	// It may take a few minutes after starting a VM for it to become available to assign to a configuration
	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdate(ctx, id, configurationAssignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return pluginsdk.RetryableError(fmt.Errorf("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("issuing creating request for %s: %+v", id, err))
		}

		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(id.ID())
	return resourceArmMaintenanceAssignmentVirtualMachineRead(d, meta)
}

func resourceArmMaintenanceAssignmentVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationassignments.ParseScopedConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("checking for presence of existing %s: %+v", *id, err)
	}

	vmId, err := virtualmachines.ParseVirtualMachineID(id.Scope)
	if err != nil {
		return err
	}

	d.Set("virtual_machine_id", vmId.ID())

	if model := resp.Model; model != nil {
		loc := location.NormalizeNilable(model.Location)
		// location isn't returned by the API
		if loc == "" {
			loc = d.Get("location").(string)
		}
		d.Set("location", loc)

		if props := model.Properties; props != nil {
			maintenanceConfigurationId := ""
			if props.MaintenanceConfigurationId != nil {
				parsedId, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationId)
				if err != nil {
					return fmt.Errorf("parsing %q: %+v", *props.MaintenanceConfigurationId, err)
				}
				maintenanceConfigurationId = parsedId.ID()
			}
			d.Set("maintenance_configuration_id", maintenanceConfigurationId)
		}
	}
	return nil
}

func resourceArmMaintenanceAssignmentVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationassignments.ParseScopedConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

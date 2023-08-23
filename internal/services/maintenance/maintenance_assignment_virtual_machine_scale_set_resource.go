// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	parseCompute "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	validateCompute "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmMaintenanceAssignmentVirtualMachineScaleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmMaintenanceAssignmentVirtualMachineScaleSetCreate,
		Read:   resourceArmMaintenanceAssignmentVirtualMachineScaleSetRead,
		Delete: resourceArmMaintenanceAssignmentVirtualMachineScaleSetDelete,

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

			if _, err := parseCompute.VirtualMachineScaleSetID(parsed.Scope); err != nil {
				return fmt.Errorf("parsing %q as a Virtual Machine Scale Set ID: %+v", parsed.Scope, err)
			}

			return nil
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AssignmentVirtualMachineScaleSetV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.Location(),

			"maintenance_configuration_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     maintenanceconfigurations.ValidateMaintenanceConfigurationID,
				DiffSuppressFunc: suppress.CaseDifference, // TODO remove in 4.0 with a work around or when https://github.com/Azure/azure-rest-api-specs/issues/8653 is fixed
			},

			"virtual_machine_scale_set_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateCompute.VirtualMachineScaleSetID,
				DiffSuppressFunc: suppress.CaseDifference, // TODO remove in 4.0
			},
		},
	}
}

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineScaleSetId, err := parseCompute.VirtualMachineScaleSetID(d.Get("virtual_machine_scale_set_id").(string))
	if err != nil {
		return err
	}

	maintenanceConfigurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Get("maintenance_configuration_id").(string))
	if err != nil {
		return err
	}

	id := configurationassignments.NewScopedConfigurationAssignmentID(virtualMachineScaleSetId.ID(), maintenanceConfigurationId.MaintenanceConfigurationName)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine_scale_set", id.ID())
	}

	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Name:     utils.String(maintenanceConfigurationId.MaintenanceConfigurationName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: utils.String(maintenanceConfigurationId.ID()),
			ResourceId:                 utils.String(virtualMachineScaleSetId.ID()),
		},
	}

	_, err = client.CreateOrUpdate(ctx, id, configurationAssignment)
	if err != nil {
		return err
	}

	d.SetId(id.ID())
	return resourceArmMaintenanceAssignmentVirtualMachineScaleSetRead(d, meta)
}

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	virtualMachineScaleSetId, err := parseCompute.VirtualMachineScaleSetID(id.Scope)
	if err != nil {
		return fmt.Errorf("parsing %q as a virtual machine scale set id: %+v", id.Scope, err)
	}

	d.Set("virtual_machine_scale_set_id", virtualMachineScaleSetId.ID())

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

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationassignments.ParseScopedConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteParent(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

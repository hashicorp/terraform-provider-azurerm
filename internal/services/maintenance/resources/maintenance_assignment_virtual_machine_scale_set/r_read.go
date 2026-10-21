// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_virtual_machine_scale_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta any) error {
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

	virtualMachineScaleSetId, err := commonids.ParseVirtualMachineScaleSetID(id.Scope)
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

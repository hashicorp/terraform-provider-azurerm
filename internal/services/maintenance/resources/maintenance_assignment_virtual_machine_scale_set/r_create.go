// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_virtual_machine_scale_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineScaleSetId, err := commonids.ParseVirtualMachineScaleSetID(d.Get("virtual_machine_scale_set_id").(string))
	if err != nil {
		return err
	}

	maintenanceConfigurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Get("maintenance_configuration_id").(string))
	if err != nil {
		return err
	}

	id := configurationassignments.NewScopedConfigurationAssignmentID(virtualMachineScaleSetId.ID(), maintenanceConfigurationId.MaintenanceConfigurationName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine_scale_set", id.ID())
		}
	}

	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Name:     pointer.To(maintenanceConfigurationId.MaintenanceConfigurationName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: pointer.To(maintenanceConfigurationId.ID()),
			ResourceId:                 pointer.To(virtualMachineScaleSetId.ID()),
		},
	}

	if _, err = client.CreateOrUpdate(ctx, id, configurationAssignment); err != nil {
		return err
	}

	d.SetId(id.ID())
	return resourceArmMaintenanceAssignmentVirtualMachineScaleSetRead(d, meta)
}

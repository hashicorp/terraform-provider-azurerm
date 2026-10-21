// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_virtual_machine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmMaintenanceAssignmentVirtualMachineCreate(d *pluginsdk.ResourceData, meta any) error {
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

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine", id.ID())
		}
	}

	// set assignment name to configuration name
	assignmentName := configurationId.MaintenanceConfigurationName
	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Name:     pointer.To(assignmentName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: pointer.To(configurationId.ID()),
			ResourceId:                 pointer.To(virtualMachineId.ID()),
		},
	}

	// It may take a few minutes after starting a VM for it to become available to assign to a configuration
	if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdate(ctx, id, configurationAssignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return pluginsdk.RetryableError(errors.New("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("issuing creating request for %s: %+v", id, err))
		}

		return nil
	}); err != nil {
		return err
	}

	d.SetId(id.ID())
	return resourceArmMaintenanceAssignmentVirtualMachineRead(d, meta)
}

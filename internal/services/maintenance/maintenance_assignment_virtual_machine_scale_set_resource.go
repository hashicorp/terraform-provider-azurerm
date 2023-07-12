// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"context"
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/parse"
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
			_, err := parse.MaintenanceAssignmentVirtualMachineScaleSetID(id)
			return err
		}),

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

	maintenanceConfigurationID, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Get("maintenance_configuration_id").(string))
	if err != nil {
		return err
	}

	configAssignmentId := configurationassignments.NewConfigurationAssignmentID(virtualMachineScaleSetId.SubscriptionId, virtualMachineScaleSetId.ResourceGroup, "Microsoft.Compute", "virtualMachineScaleSets", virtualMachineScaleSetId.Name, maintenanceConfigurationID.MaintenanceConfigurationName)

	existingList, err := getMaintenanceAssignmentVirtualMachineScaleSet(ctx, client, virtualMachineScaleSetId)
	if err != nil {
		return err
	}
	if existingList != nil && len(*existingList) > 0 {
		existing := (*existingList)[0]
		if existing.Id != nil && *existing.Id != "" {
			return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine_scale_set", configAssignmentId.ID())
		}
	}

	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Name:     utils.String(maintenanceConfigurationID.MaintenanceConfigurationName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: utils.String(maintenanceConfigurationID.ID()),
			ResourceId:                 utils.String(virtualMachineScaleSetId.ID()),
		},
	}

	_, err = client.CreateOrUpdate(ctx, configAssignmentId, configurationAssignment)
	if err != nil {
		return err
	}

	d.SetId(configAssignmentId.ID())
	return resourceArmMaintenanceAssignmentVirtualMachineScaleSetRead(d, meta)
}

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := getMaintenanceAssignmentVirtualMachineScaleSet(ctx, client, id.VirtualMachineScaleSetId)
	if err != nil {
		return err
	}
	if resp == nil || len(*resp) == 0 {
		d.SetId("")
		return nil
	}
	assignment := (*resp)[0]
	if assignment.Id == nil || *assignment.Id == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (virtual machine scale set ID id: %q", id.VirtualMachineScaleSetIdRaw)
	}

	// in list api, `ResourceID` returned is always nil
	virtualMachineScaleSetId := ""
	if id.VirtualMachineScaleSetId != nil {
		virtualMachineScaleSetId = id.VirtualMachineScaleSetId.ID()
	}
	d.Set("virtual_machine_scale_set_id", virtualMachineScaleSetId)
	if props := assignment.Properties; props != nil {
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
	return nil
}

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	maintenanceAssignmentVmScaleSetId, err := parse.MaintenanceAssignmentVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	id := configurationassignments.NewConfigurationAssignmentID(maintenanceAssignmentVmScaleSetId.VirtualMachineScaleSetId.SubscriptionId, maintenanceAssignmentVmScaleSetId.VirtualMachineScaleSetId.ResourceGroup, "Microsoft.Compute", "virtualMachineScaleSets", maintenanceAssignmentVmScaleSetId.VirtualMachineScaleSetId.Name, maintenanceAssignmentVmScaleSetId.Name)

	if _, err := client.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", maintenanceAssignmentVmScaleSetId.VirtualMachineScaleSetIdRaw, err)
	}

	return nil
}

func getMaintenanceAssignmentVirtualMachineScaleSet(ctx context.Context, client *configurationassignments.ConfigurationAssignmentsClient, vmScaleSetId *parseCompute.VirtualMachineScaleSetId) (result *[]configurationassignments.ConfigurationAssignment, err error) {
	id := configurationassignments.NewProviderID(vmScaleSetId.SubscriptionId, vmScaleSetId.ResourceGroup, "Microsoft.Compute", "virtualMachineScaleSets", vmScaleSetId.Name)

	resp, err := client.List(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			err = fmt.Errorf("checking for presence of existing Maintenance assignment (virtual machine scale set ID: %q): %+v", vmScaleSetId.ID(), err)
			return
		}
	}
	return resp.Model.Value, nil
}

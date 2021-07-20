package maintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/maintenance/mgmt/2021-05-01/maintenance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	parseCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	validateCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			"location": azure.SchemaLocation(),

			"maintenance_configuration_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MaintenanceConfigurationID,
			},

			"virtual_machine_scale_set_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateCompute.VirtualMachineScaleSetID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineScaleSetIdRaw := d.Get("virtual_machine_scale_set_id").(string)
	virtualMachineScaleSetId, _ := parseCompute.VirtualMachineScaleSetID(virtualMachineScaleSetIdRaw)

	existingList, err := getMaintenanceAssignmentVirtualMachineScaleSet(ctx, client, virtualMachineScaleSetId, virtualMachineScaleSetIdRaw)
	if err != nil {
		return err
	}
	if existingList != nil && len(*existingList) > 0 {
		existing := (*existingList)[0]
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine_scale_set", *existing.ID)
		}
	}

	maintenanceConfigurationID := d.Get("maintenance_configuration_id").(string)
	configurationId, _ := parse.MaintenanceConfigurationIDInsensitively(maintenanceConfigurationID)

	// set assignment name to configuration name
	assignmentName := configurationId.Name
	configurationAssignment := maintenance.ConfigurationAssignment{
		Name:     utils.String(assignmentName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		ConfigurationAssignmentProperties: &maintenance.ConfigurationAssignmentProperties{
			MaintenanceConfigurationID: utils.String(maintenanceConfigurationID),
			ResourceID:                 utils.String(virtualMachineScaleSetIdRaw),
		},
	}

	_, err = client.CreateOrUpdate(ctx, virtualMachineScaleSetId.ResourceGroup, "Microsoft.Compute", "virtualMachineScaleSets", virtualMachineScaleSetId.Name, assignmentName, configurationAssignment)
	if err != nil {
		return err
	}

	resp, err := getMaintenanceAssignmentVirtualMachineScaleSet(ctx, client, virtualMachineScaleSetId, virtualMachineScaleSetIdRaw)
	if err != nil {
		return err
	}
	if resp == nil || len(*resp) == 0 {
		return fmt.Errorf("could not find Maintenance assignment (virtual machine scale set ID: %q)", virtualMachineScaleSetIdRaw)
	}
	assignment := (*resp)[0]
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (virtual machine scale set ID %q)", virtualMachineScaleSetIdRaw)
	}

	d.SetId(*assignment.ID)
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

	resp, err := getMaintenanceAssignmentVirtualMachineScaleSet(ctx, client, id.VirtualMachineScaleSetId, id.VirtualMachineScaleSetIdRaw)
	if err != nil {
		return err
	}
	if resp == nil || len(*resp) == 0 {
		d.SetId("")
		return nil
	}
	assignment := (*resp)[0]
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (virtual machine scale set ID id: %q", id.VirtualMachineScaleSetIdRaw)
	}

	// in list api, `ResourceID` returned is always nil
	virtualMachineScaleSetId := ""
	if id.VirtualMachineScaleSetId != nil {
		virtualMachineScaleSetId = id.VirtualMachineScaleSetId.ID()
	}
	d.Set("virtual_machine_scale_set_id", virtualMachineScaleSetId)
	if props := assignment.ConfigurationAssignmentProperties; props != nil {
		d.Set("maintenance_configuration_id", props.MaintenanceConfigurationID)
	}
	return nil
}

func resourceArmMaintenanceAssignmentVirtualMachineScaleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.VirtualMachineScaleSetId.ResourceGroup, "Microsoft.Compute", "virtualMachineScaleSets", id.VirtualMachineScaleSetId.Name, id.Name); err != nil {
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", id.VirtualMachineScaleSetIdRaw, err)
	}

	return nil
}

func getMaintenanceAssignmentVirtualMachineScaleSet(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id *parseCompute.VirtualMachineScaleSetId, virtualMachineScaleSetId string) (result *[]maintenance.ConfigurationAssignment, err error) {
	resp, err := client.List(ctx, id.ResourceGroup, "Microsoft.Compute", "virtualMachineScaleSets", id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			err = fmt.Errorf("checking for presence of existing Maintenance assignment (virtual machine scale set ID: %q): %+v", virtualMachineScaleSetId, err)
			return
		}
	}
	return resp.Value, nil
}

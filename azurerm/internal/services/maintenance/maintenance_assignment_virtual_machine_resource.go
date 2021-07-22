package maintenance

import (
	"context"
	"fmt"
	"strings"
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
			_, err := parse.MaintenanceAssignmentVirtualMachineID(id)
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

			"virtual_machine_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateCompute.VirtualMachineID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentVirtualMachineCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineIdRaw := d.Get("virtual_machine_id").(string)
	virtualMachineId, _ := parseCompute.VirtualMachineID(virtualMachineIdRaw)

	existingList, err := getMaintenanceAssignmentVirtualMachine(ctx, client, virtualMachineId, virtualMachineIdRaw)
	if err != nil {
		return err
	}
	if existingList != nil && len(*existingList) > 0 {
		existing := (*existingList)[0]
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine", *existing.ID)
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
			ResourceID:                 utils.String(virtualMachineIdRaw),
		},
	}

	// It may take a few minutes after starting a VM for it to become available to assign to a configuration
	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdate(ctx, virtualMachineId.ResourceGroup, "Microsoft.Compute", "virtualMachines", virtualMachineId.Name, assignmentName, configurationAssignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return pluginsdk.RetryableError(fmt.Errorf("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("issuing creating request for Maintenance Assignment (virtual machine ID %q): %+v", virtualMachineIdRaw, err))
		}

		return nil
	}) //lintignore:R006
	if err != nil {
		return err
	}

	resp, err := getMaintenanceAssignmentVirtualMachine(ctx, client, virtualMachineId, virtualMachineIdRaw)
	if err != nil {
		return err
	}
	if resp == nil || len(*resp) == 0 {
		return fmt.Errorf("could not find Maintenance assignment (virtual machine ID: %q)", virtualMachineIdRaw)
	}
	assignment := (*resp)[0]
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (virtual machine ID %q)", virtualMachineIdRaw)
	}

	d.SetId(*assignment.ID)
	return resourceArmMaintenanceAssignmentVirtualMachineRead(d, meta)
}

func resourceArmMaintenanceAssignmentVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := getMaintenanceAssignmentVirtualMachine(ctx, client, id.VirtualMachineId, id.VirtualMachineIdRaw)
	if err != nil {
		return err
	}
	if resp == nil || len(*resp) == 0 {
		d.SetId("")
		return nil
	}
	assignment := (*resp)[0]
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (virtual machine ID id: %q", id.VirtualMachineIdRaw)
	}

	// in list api, `ResourceID` returned is always nil
	virtualMachineId := ""
	if id.VirtualMachineId != nil {
		virtualMachineId = id.VirtualMachineId.ID()
	}
	d.Set("virtual_machine_id", virtualMachineId)
	if props := assignment.ConfigurationAssignmentProperties; props != nil {
		d.Set("maintenance_configuration_id", props.MaintenanceConfigurationID)
	}
	return nil
}

func resourceArmMaintenanceAssignmentVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.VirtualMachineId.ResourceGroup, "Microsoft.Compute", "virtualMachines", id.VirtualMachineId.Name, id.Name); err != nil {
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", id.VirtualMachineIdRaw, err)
	}

	return nil
}

func getMaintenanceAssignmentVirtualMachine(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id *parseCompute.VirtualMachineId, virtualMachineId string) (result *[]maintenance.ConfigurationAssignment, err error) {
	resp, err := client.List(ctx, id.ResourceGroup, "Microsoft.Compute", "virtualMachines", id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			err = fmt.Errorf("checking for presence of existing Maintenance assignment (virtual machine ID: %q): %+v", virtualMachineId, err)
			return
		}
	}
	return resp.Value, nil
}

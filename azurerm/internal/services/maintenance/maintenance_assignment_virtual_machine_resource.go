package maintenance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/maintenance/mgmt/2018-06-01-preview/maintenance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	parseCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	validateCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMaintenanceAssignmentVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMaintenanceAssignmentVirtualMachineCreate,
		Read:   resourceArmMaintenanceAssignmentVirtualMachineRead,
		Delete: resourceArmMaintenanceAssignmentVirtualMachineDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MaintenanceAssignmentVirtualMachineID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"location": azure.SchemaLocation(),

			"maintenance_configuration_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MaintenanceConfigurationID,
			},

			"virtual_machine_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateCompute.VirtualMachineID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineIdRaw := d.Get("virtual_machine_id").(string)
	virtualMachineId, _ := parseCompute.VirtualMachineID(virtualMachineIdRaw)

	existing, err := getMaintenanceAssignmentVirtualMachine(ctx, client, virtualMachineId, virtualMachineIdRaw)
	if err != nil {
		return err
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine", *existing.ID)
	}

	maintenanceConfigurationID := d.Get("maintenance_configuration_id").(string)
	configurationId, _ := parse.MaintenanceConfigurationID(maintenanceConfigurationID)

	// set assignment name to configuration name
	assignmentName := configurationId.Name
	assignment := maintenance.ConfigurationAssignment{
		Name:     utils.String(assignmentName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		ConfigurationAssignmentProperties: &maintenance.ConfigurationAssignmentProperties{
			MaintenanceConfigurationID: utils.String(maintenanceConfigurationID),
			ResourceID:                 utils.String(virtualMachineIdRaw),
		},
	}

	// It may take a few minutes after starting a VM for it to become available to assign to a configuration
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if _, err := client.CreateOrUpdate(ctx, virtualMachineId.ResourceGroup, "Microsoft.Compute", "virtualMachines", virtualMachineId.Name, assignmentName, assignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return resource.RetryableError(fmt.Errorf("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return resource.NonRetryableError(fmt.Errorf("issuing creating request for Maintenance Assignment (virtual machine ID %q): %+v", virtualMachineIdRaw, err))
		}

		return nil
	})
	if err != nil {
		return err
	}

	resp, err := getMaintenanceAssignmentVirtualMachine(ctx, client, virtualMachineId, virtualMachineIdRaw)
	if err != nil {
		return err
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (virtual machine ID %q)", virtualMachineIdRaw)
	}

	d.SetId(*resp.ID)
	return resourceArmMaintenanceAssignmentVirtualMachineRead(d, meta)
}

func resourceArmMaintenanceAssignmentVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	assignment, err := getMaintenanceAssignmentVirtualMachine(ctx, client, id.VirtualMachineId, id.VirtualMachineIdRaw)
	if err != nil {
		return err
	}
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (virtual machine ID id: %q", id.VirtualMachineIdRaw)
	}

	// in list api, `ResourceID` returned is always nil
	d.Set("virtual_machine_id", id.VirtualMachineIdRaw)
	if props := assignment.ConfigurationAssignmentProperties; props != nil {
		d.Set("maintenance_configuration_id", props.MaintenanceConfigurationID)
	}
	return nil
}

func resourceArmMaintenanceAssignmentVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
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

func getMaintenanceAssignmentVirtualMachine(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id *parseCompute.VirtualMachineId, virtualMachineId string) (result maintenance.ConfigurationAssignment, err error) {
	resp, err := client.List(ctx, id.ResourceGroup, "Microsoft.Compute", "virtualMachines", id.Name)

	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			err = fmt.Errorf("checking for presence of existing Maintenance assignment (virtual machine ID: %q): %+v", virtualMachineId, err)
			return
		}
		return result, nil
	}
	if resp.Value == nil || len(*resp.Value) == 0 {
		err = fmt.Errorf("could not find Maintenance assignment (virtual machine ID: %q)", virtualMachineId)
		return
	}

	return (*resp.Value)[0], nil
}

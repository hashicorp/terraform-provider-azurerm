package maintenance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2021-05-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2021-05-01/maintenanceconfigurations"
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
			"location": commonschema.Location(),

			"maintenance_configuration_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     maintenanceconfigurations.ValidateMaintenanceConfigurationID,
				DiffSuppressFunc: suppress.CaseDifference, // TODO remove in 4.0 with a work around or when https://github.com/Azure/azure-rest-api-specs/issues/8653 is fixed
			},

			"virtual_machine_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateCompute.VirtualMachineID,
				DiffSuppressFunc: suppress.CaseDifference, // TODO remove in 4.0
			},
		},
	}
}

func resourceArmMaintenanceAssignmentVirtualMachineCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineId, err := parseCompute.VirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}

	configurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Get("maintenance_configuration_id").(string))
	if err != nil {
		return err
	}

	existingList, err := getMaintenanceAssignmentVirtualMachine(ctx, client, virtualMachineId, virtualMachineId.ID())
	if err != nil {
		return err
	}
	if existingList != nil && len(*existingList) > 0 {
		existing := (*existingList)[0]
		if existing.Id != nil && *existing.Id != "" {
			return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine", configurationId.ID())
		}
	}

	// set assignment name to configuration name
	assignmentName := configurationId.ResourceName
	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Name:     utils.String(assignmentName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: utils.String(configurationId.ID()),
			ResourceId:                 utils.String(virtualMachineId.ID()),
		},
	}

	id := configurationassignments.NewConfigurationAssignmentID(virtualMachineId.SubscriptionId, virtualMachineId.ResourceGroup, "Microsoft.Compute", "virtualMachines", virtualMachineId.Name, assignmentName)

	// It may take a few minutes after starting a VM for it to become available to assign to a configuration
	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdate(ctx, id, configurationAssignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return pluginsdk.RetryableError(fmt.Errorf("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("issuing creating request for Maintenance Assignment (virtual machine ID %q): %+v", virtualMachineId.ID(), err))
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
	if assignment.Id == nil || *assignment.Id == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (virtual machine ID id: %q", id.VirtualMachineIdRaw)
	}

	// in list api, `ResourceID` returned is always nil
	virtualMachineId := ""
	if id.VirtualMachineId != nil {
		virtualMachineId = id.VirtualMachineId.ID()
	}
	d.Set("virtual_machine_id", virtualMachineId)
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

func resourceArmMaintenanceAssignmentVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId, err := parse.MaintenanceAssignmentVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	id := configurationassignments.NewConfigurationAssignmentID(vmId.VirtualMachineId.SubscriptionId, vmId.VirtualMachineId.ResourceGroup, "Microsoft.Compute", "virtualMachines", vmId.VirtualMachineId.Name, vmId.Name)

	if _, err := client.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", vmId.VirtualMachineIdRaw, err)
	}

	return nil
}

func getMaintenanceAssignmentVirtualMachine(ctx context.Context, client *configurationassignments.ConfigurationAssignmentsClient, vmId *parseCompute.VirtualMachineId, virtualMachineId string) (result *[]configurationassignments.ConfigurationAssignment, err error) {

	id := configurationassignments.NewProviderID(vmId.SubscriptionId, vmId.ResourceGroup, "Microsoft.Compute", "virtualMachines", vmId.Name)
	resp, err := client.List(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			err = fmt.Errorf("checking for presence of existing Maintenance assignment (virtual machine ID: %q): %+v", virtualMachineId, err)
			return
		}
	}

	return resp.Model.Value, nil
}

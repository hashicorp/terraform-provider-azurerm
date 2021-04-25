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

func resourceArmMaintenanceAssignmentDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMaintenanceAssignmentDedicatedHostCreate,
		Read:   resourceArmMaintenanceAssignmentDedicatedHostRead,
		Delete: resourceArmMaintenanceAssignmentDedicatedHostDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.MaintenanceAssignmentDedicatedHostID(id)
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

			"dedicated_host_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateCompute.DedicatedHostID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dedicatedHostIdRaw := d.Get("dedicated_host_id").(string)
	dedicatedHostId, _ := parseCompute.DedicatedHostID(dedicatedHostIdRaw)

	existing, err := getMaintenanceAssignmentDedicatedHost(ctx, client, dedicatedHostId, dedicatedHostIdRaw)
	if err != nil {
		return err
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_maintenance_assignment_dedicated_host", *existing.ID)
	}

	maintenanceConfigurationID := d.Get("maintenance_configuration_id").(string)
	configurationId, _ := parse.MaintenanceConfigurationIDInsensitively(maintenanceConfigurationID)

	// set assignment name to configuration name
	assignmentName := configurationId.Name
	assignment := maintenance.ConfigurationAssignment{
		Name:     utils.String(assignmentName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		ConfigurationAssignmentProperties: &maintenance.ConfigurationAssignmentProperties{
			MaintenanceConfigurationID: utils.String(maintenanceConfigurationID),
			ResourceID:                 utils.String(dedicatedHostIdRaw),
		},
	}

	// It may take a few minutes after starting a VM for it to become available to assign to a configuration
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if _, err := client.CreateOrUpdateParent(ctx, dedicatedHostId.ResourceGroup, "Microsoft.Compute", "hostGroups", dedicatedHostId.HostGroupName, "hosts", dedicatedHostId.HostName, assignmentName, assignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return resource.RetryableError(fmt.Errorf("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return resource.NonRetryableError(fmt.Errorf("issuing creating request for Maintenance Assignment (Dedicated Host ID %q): %+v", dedicatedHostIdRaw, err))
		}

		return nil
	})
	if err != nil {
		return err
	}

	resp, err := getMaintenanceAssignmentDedicatedHost(ctx, client, dedicatedHostId, dedicatedHostIdRaw)
	if err != nil {
		return err
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (Dedicated Host ID %q)", dedicatedHostIdRaw)
	}

	d.SetId(*resp.ID)
	return resourceArmMaintenanceAssignmentDedicatedHostRead(d, meta)
}

func resourceArmMaintenanceAssignmentDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentDedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	assignment, err := getMaintenanceAssignmentDedicatedHost(ctx, client, id.DedicatedHostId, id.DedicatedHostIdRaw)
	if err != nil {
		return err
	}
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (Dedicated Host ID: %q", id.DedicatedHostIdRaw)
	}

	dedicatedHostId := ""
	if id.DedicatedHostId != nil {
		dedicatedHostId = id.DedicatedHostId.ID()
	}
	d.Set("dedicated_host_id", dedicatedHostId)

	if props := assignment.ConfigurationAssignmentProperties; props != nil {
		d.Set("maintenance_configuration_id", props.MaintenanceConfigurationID)
	}
	return nil
}

func resourceArmMaintenanceAssignmentDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentDedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteParent(ctx, id.DedicatedHostId.ResourceGroup, "Microsoft.Compute", "hostGroups", id.DedicatedHostId.HostGroupName, "hosts", id.DedicatedHostId.HostName, id.Name); err != nil {
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", id.DedicatedHostIdRaw, err)
	}

	return nil
}

func getMaintenanceAssignmentDedicatedHost(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id *parseCompute.DedicatedHostId, dedicatedHostId string) (result maintenance.ConfigurationAssignment, err error) {
	resp, err := client.ListParent(ctx, id.ResourceGroup, "Microsoft.Compute", "hostGroups", id.HostGroupName, "hosts", id.HostName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			err = fmt.Errorf("checking for presence of existing Maintenance assignment (Dedicated Host ID %q): %+v", dedicatedHostId, err)
			return
		}
		return result, nil
	}
	if resp.Value == nil || len(*resp.Value) == 0 {
		err = fmt.Errorf("could not find Maintenance assignment (Dedicated Host ID %q)", dedicatedHostId)
		return
	}

	return (*resp.Value)[0], nil
}

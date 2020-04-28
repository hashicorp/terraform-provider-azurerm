package maintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/maintenance/mgmt/2018-06-01-preview/maintenance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	validateCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMaintenanceAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMaintenanceAssignmentCreate,
		Read:   resourceArmMaintenanceAssignmentRead,
		Update: nil,
		Delete: resourceArmMaintenanceAssignmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MaintenanceAssignmentID(id)
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

			"target_resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validateCompute.DedicatedHostID,
					validateCompute.VirtualMachineID,
				),
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	targetResourceId := d.Get("target_resource_id").(string)
	id, _ := parse.TargetResourceID(targetResourceId)

	existing, err := getMaintenanceAssignment(ctx, client, id, targetResourceId)
	if err != nil {
		return err
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_maintenance_assignment", *existing.ID)
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
			ResourceID:                 utils.String(targetResourceId),
		},
	}

	if _, err := createMaintenanceAssignment(ctx, client, id, assignmentName, &assignment); err != nil {
		return fmt.Errorf("creating Maintenance Assignment (target resource id: %q): %+v", targetResourceId, err)
	}

	resp, err := getMaintenanceAssignment(ctx, client, id, targetResourceId)
	if err != nil {
		return err
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (target resource id %q)", targetResourceId)
	}

	d.SetId(*resp.ID)
	return resourceArmMaintenanceAssignmentRead(d, meta)
}

func resourceArmMaintenanceAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentID(d.Id())
	if err != nil {
		return err
	}

	assignment, err := getMaintenanceAssignment(ctx, client, id.TargetResourceId, id.ResourceId)
	if err != nil {
		return err
	}
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (target resource id: %q", id.ResourceId)
	}

	// in list api, `ResourceID` returned is always nil
	d.Set("target_resource_id", id.ResourceId)
	if props := assignment.ConfigurationAssignmentProperties; props != nil {
		d.Set("maintenance_configuration_id", props.MaintenanceConfigurationID)
	}
	return nil
}

func resourceArmMaintenanceAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := deleteMaintenanceAssignment(ctx, client, id); err != nil {
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", id.ResourceId, err)
	}

	return nil
}

func getMaintenanceAssignment(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id *parse.TargetResourceId, targetResourceId string) (maintenance.ConfigurationAssignment, error) {
	var listResp maintenance.ListConfigurationAssignmentsResult
	var err error
	if id.HasParentResource {
		listResp, err = client.ListParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName)
	} else {
		listResp, err = client.List(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName)
	}
	if err != nil {
		if !utils.ResponseWasNotFound(listResp.Response) {
			return maintenance.ConfigurationAssignment{}, fmt.Errorf("checking for presence of existing Maintenance assignment (target resource id %q): %+v", targetResourceId, err)
		}
		return maintenance.ConfigurationAssignment{}, nil
	}
	if listResp.Value == nil || len(*listResp.Value) == 0 {
		return maintenance.ConfigurationAssignment{}, fmt.Errorf("could not find Maintenance assignment (target resource id %q)", targetResourceId)
	}

	return (*listResp.Value)[0], nil
}

func createMaintenanceAssignment(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id *parse.TargetResourceId, assignmentName string, assignment *maintenance.ConfigurationAssignment) (maintenance.ConfigurationAssignment, error) {
	if id.HasParentResource {
		return client.CreateOrUpdateParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName, assignmentName, *assignment)
	} else {
		return client.CreateOrUpdate(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName, assignmentName, *assignment)
	}
}

func deleteMaintenanceAssignment(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id *parse.MaintenanceAssignmentId) (maintenance.ConfigurationAssignment, error) {
	if id.HasParentResource {
		return client.DeleteParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName, id.Name)
	} else {
		return client.Delete(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName, id.Name)
	}
}

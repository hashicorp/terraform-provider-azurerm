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

	id, _ := parse.TargetResourceID(d.Get("target_resource_id").(string))

	existing, err := getMaintenanceAssignment(ctx, client, id)
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
			ResourceID:                 utils.String(id.ID()),
		},
	}

	if _, err := createMaintenanceAssignment(ctx, client, id, assignmentName, &assignment); err != nil {
		return fmt.Errorf("creating Maintenance Assignment (target resource id: %q): %+v", id.ID(), err)
	}

	resp, err := getMaintenanceAssignment(ctx, client, id)
	if err != nil {
		return err
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (target resource id %q)", id.ID())
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

	assignment, err := getMaintenanceAssignment(ctx, client, id.TargetResourceId)
	if err != nil {
		return err
	}
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (target resource id: %q", id.TargetResourceId.ID())
	}

	// in list api, `ResourceID` returned is always nil
	d.Set("target_resource_id", id.TargetResourceId.ID())
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
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", id.TargetResourceId.ID(), err)
	}

	return nil
}

func getMaintenanceAssignment(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id parse.TargetResourceId) (result maintenance.ConfigurationAssignment, err error) {
	var listResp maintenance.ListConfigurationAssignmentsResult

	switch v := id.(type) {
	case parse.ScopeResource:
		listResp, err = client.List(ctx, v.ResourceGroup, v.ResourceProvider, v.ResourceType, v.ResourceName)
	case parse.ScopeInResource:
		listResp, err = client.ListParent(ctx, v.ResourceGroup, v.ResourceProvider, v.ResourceParentType, v.ResourceParentName, v.ResourceType, v.ResourceName)
	default:
		err = fmt.Errorf("wrong type of target resource id: %+v", id)
		return
	}

	if err != nil {
		if !utils.ResponseWasNotFound(listResp.Response) {
			err = fmt.Errorf("checking for presence of existing Maintenance assignment (target resource id %q): %+v", id.ID(), err)
			return
		}
		return result, nil
	}
	if listResp.Value == nil || len(*listResp.Value) == 0 {
		err = fmt.Errorf("could not find Maintenance assignment (target resource id %q)", id.ID())
		return
	}

	return (*listResp.Value)[0], nil
}

func createMaintenanceAssignment(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id parse.TargetResourceId, assignmentName string, assignment *maintenance.ConfigurationAssignment) (result maintenance.ConfigurationAssignment, err error) {
	switch v := id.(type) {
	case parse.ScopeResource:
		return client.CreateOrUpdate(ctx, v.ResourceGroup, v.ResourceProvider, v.ResourceType, v.ResourceName, assignmentName, *assignment)
	case parse.ScopeInResource:
		return client.CreateOrUpdateParent(ctx, v.ResourceGroup, v.ResourceProvider, v.ResourceParentType, v.ResourceParentName, v.ResourceType, v.ResourceName, assignmentName, *assignment)
	default:
		err = fmt.Errorf("wrong type of target resource id: %+v", id)
		return
	}
}

func deleteMaintenanceAssignment(ctx context.Context, client *maintenance.ConfigurationAssignmentsClient, id *parse.MaintenanceAssignmentId) (result maintenance.ConfigurationAssignment, err error) {
	switch v := id.TargetResourceId.(type) {
	case parse.ScopeResource:
		return client.Delete(ctx, v.ResourceGroup, v.ResourceProvider, v.ResourceType, v.ResourceName, id.Name)
	case parse.ScopeInResource:
		return client.DeleteParent(ctx, v.ResourceGroup, v.ResourceProvider, v.ResourceParentType, v.ResourceParentName, v.ResourceType, v.ResourceName, id.Name)
	default:
		err = fmt.Errorf("wrong type of target resource id: %+v", id)
		return
	}
}

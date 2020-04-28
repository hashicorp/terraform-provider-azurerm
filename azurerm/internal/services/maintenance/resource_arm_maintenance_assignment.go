package maintenance

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/maintenance/mgmt/2018-06-01-preview/maintenance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	maintenanceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/validate"
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
				ValidateFunc: maintenanceValidate.MaintenanceConfigurationID,
			},

			"target_resource_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	targetResourceId := d.Get("target_resource_id").(string)
	id, err := parse.TargetResourceID(targetResourceId)
	if err != nil {
		return err
	}

	var listResp maintenance.ListConfigurationAssignmentsResult
	if id.HasParentResource {
		listResp, err = client.ListParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName)
	} else {
		listResp, err = client.List(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName)
	}
	if err != nil {
		if !utils.ResponseWasNotFound(listResp.Response) {
			return fmt.Errorf("checking for presense of existing Maintenance assignment to resource %q: %+v", targetResourceId, err)
		}
	}
	if listResp.Value != nil && len(*listResp.Value) > 0 {
		return tf.ImportAsExistsError("azurerm_maintenance_assignment", *(*listResp.Value)[0].ID)
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

	var resp maintenance.ConfigurationAssignment
	if id.HasParentResource {
		resp, err = client.CreateOrUpdateParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName, assignmentName, assignment)
	} else {
		resp, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName, assignmentName, assignment)
	}
	if err != nil {
		return fmt.Errorf("creating Maintenance Assignment to resource %q: %+v", targetResourceId, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Maintenance Assignment to resource %q", targetResourceId)
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

	var listResp maintenance.ListConfigurationAssignmentsResult
	if id.HasParentResource {
		listResp, err = client.ListParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName)
	} else {
		listResp, err = client.List(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName)
	}
	if err != nil {
		if !utils.ResponseWasNotFound(listResp.Response) {
			return fmt.Errorf("checking for present of existing Maintenance assignment to resource %q: %+v", id.ResourceId, err)
		}
		return fmt.Errorf("listing Maintenance assignment to resource %q: %+v", id.ResourceId, err)
	}
	if listResp.Value == nil || len(*listResp.Value) == 0 {
		return fmt.Errorf("could not find Maintenance assignment to resource %q", id.ResourceId)
	}

	assignment := (*listResp.Value)[0]
	if assignment.ID == nil || *assignment.ID == "" {
		return fmt.Errorf("cannot read Maintenance Assignment to resource %q ID", id.ResourceId)
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

	if id.HasParentResource {
		_, err = client.DeleteParent(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName, id.Name)
	} else {
		_, err = client.Delete(ctx, id.ResourceGroup, id.ResourceProvider, id.ResourceType, id.ResourceName, id.Name)
	}
	if err != nil {
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", id.ResourceId, err)
	}

	return nil
}

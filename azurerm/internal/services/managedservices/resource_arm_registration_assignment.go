package managedservices

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRegistrationAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRegistrationAssignmentCreateUpdate,
		Read:   resourceArmRegistrationAssignmentRead,
		Update: resourceArmRegistrationAssignmentCreateUpdate,
		Delete: resourceArmRegistrationAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"registration_assignment_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"registration_definition_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmRegistrationAssignmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.RegistrationAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	registrationAssignmentId := d.Get("registration_assignment_id").(string)
	if registrationAssignmentId == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Registration Assignment: %+v", err)
		}

		registrationAssignmentId = uuid
	}

	scope := d.Get("scope").(string)
	expandRegistrationDefinition := true

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, scope, registrationAssignmentId, &expandRegistrationDefinition)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Registration Assignment %q (Scope %q): %+v", registrationAssignmentId, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_registration_assignment", *existing.ID)
		}
	}

	parameters := managedservices.RegistrationAssignment{
		Properties: &managedservices.RegistrationAssignmentProperties{
			RegistrationDefinitionID: utils.String(d.Get("registration_definition_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, scope, registrationAssignmentId, parameters); err != nil {
		return fmt.Errorf("Error Creating/Updating Registration Assignment %q (Scope %q): %+v", registrationAssignmentId, scope, err)
	}

	read, err := client.Get(ctx, scope, registrationAssignmentId, &expandRegistrationDefinition)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Registration Assignment %q ID (scope %q) ID", registrationAssignmentId, scope)
	}

	d.SetId(*read.ID)

	return resourceArmRegistrationAssignmentRead(d, meta)
}

func resourceArmRegistrationAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.RegistrationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureRegistrationAssignmentId(d.Id())
	if err != nil {
		return err
	}
	expandRegistrationDefinition := true

	resp, err := client.Get(ctx, id.scope, id.registrationAssignmentId, &expandRegistrationDefinition)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Registration Assignment '%s' was not found (Scope '%s')", id.registrationAssignmentId, id.scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Registration Assignment %q (Scope %q): %+v", id.registrationAssignmentId, id.scope, err)
	}

	d.Set("registration_assignment_id", resp.Name)
	d.Set("scope", id.scope)

	if props := resp.Properties; props != nil {
		d.Set("registration_definition_id", props.RegistrationDefinitionID)
	}

	return nil
}

type registrationAssignmentID struct {
	scope                    string
	registrationAssignmentId string
}

func parseAzureRegistrationAssignmentId(id string) (*registrationAssignmentID, error) {
	segments := strings.Split(id, "/providers/Microsoft.ManagedServices/registrationAssignments/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.ManagedServices/registrationAssignments/{name} - got %d segments", len(segments))
	}

	azureRegistrationAssignmentId := registrationAssignmentID{
		scope:                    segments[0],
		registrationAssignmentId: segments[1],
	}

	return &azureRegistrationAssignmentId, nil
}

func resourceArmRegistrationAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.RegistrationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureRegistrationAssignmentId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.scope, id.registrationAssignmentId)
	if err != nil {
		return fmt.Errorf("Error deleting Registration Assignment %q at Scope %q: %+v", id.registrationAssignmentId, id.scope, err)
	}

	// The sleep is needed to ensure the registration assignment is successfully deleted.
	// Bug # is logged with the Product team to track this issue.
	time.Sleep(20 * time.Second)

	return nil
}

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

func resourceArmLighthouseAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLighthouseAssignmentCreateUpdate,
		Read:   resourceArmLighthouseAssignmentRead,
		Update: resourceArmLighthouseAssignmentCreateUpdate,
		Delete: resourceArmLighthouseAssignmentDelete,
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
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"registration_definition_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"expand_registration_definition": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmLighthouseAssignmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.LighthouseAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseAssignmentID := d.Get("registration_assignment_id").(string)
	if lighthouseAssignmentID == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Lighthouse Assignment: %+v", err)
		}

		lighthouseAssignmentID = uuid
	}

	scope := d.Get("scope").(string)
	expandLighthouseDefinition := d.Get("expand_registration_definition").(bool)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, scope, lighthouseAssignmentID, &expandLighthouseDefinition)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Lighthouse Assignment %q (Scope %q): %+v", lighthouseAssignmentID, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_lighthouse_assignment", *existing.ID)
		}
	}

	parameters := managedservices.RegistrationAssignment{
		Properties: &managedservices.RegistrationAssignmentProperties{
			RegistrationDefinitionID: utils.String(d.Get("registration_definition_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, scope, lighthouseAssignmentID, parameters); err != nil {
		return fmt.Errorf("Error Creating/Updating Lighthouse Assignment %q (Scope %q): %+v", lighthouseAssignmentID, scope, err)
	}

	read, err := client.Get(ctx, scope, lighthouseAssignmentID, &expandLighthouseDefinition)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Lighthouse Assignment %q ID (scope %q) ID", lighthouseAssignmentID, scope)
	}

	d.SetId(*read.ID)

	return resourceArmLighthouseAssignmentRead(d, meta)
}

func resourceArmLighthouseAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.LighthouseAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureLighthouseAssignmentID(d.Id())
	if err != nil {
		return err
	}
	expandLighthouseDefinition := d.Get("expand_registration_definition").(bool)

	resp, err := client.Get(ctx, id.scope, id.lighthouseAssignmentID, &expandLighthouseDefinition)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Lighthouse Assignment '%s' was not found (Scope '%s')", id.lighthouseAssignmentID, id.scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Lighthouse Assignment %q (Scope %q): %+v", id.lighthouseAssignmentID, id.scope, err)
	}

	d.Set("registration_assignment_id", resp.Name)
	d.Set("scope", id.scope)

	if props := resp.Properties; props != nil {
		d.Set("registration_definition_id", props.RegistrationDefinitionID)
	}

	return nil
}

type lighthouseAssignmentID struct {
	scope                  string
	lighthouseAssignmentID string
}

func parseAzureLighthouseAssignmentID(id string) (*lighthouseAssignmentID, error) {
	segments := strings.Split(id, "/providers/Microsoft.ManagedServices/registrationAssignments/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.ManagedServices/registrationAssignments/{name} - got %d segments", len(segments))
	}

	azureLighthouseAssignmentID := lighthouseAssignmentID{
		scope:                  segments[0],
		lighthouseAssignmentID: segments[1],
	}

	return &azureLighthouseAssignmentID, nil
}

func resourceArmLighthouseAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.LighthouseAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureLighthouseAssignmentID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.scope, id.lighthouseAssignmentID)
	if err != nil {
		return fmt.Errorf("Error deleting Lighthouse Assignment %q at Scope %q: %+v", id.lighthouseAssignmentID, id.scope, err)
	}

	// The sleep is needed to ensure the lighthouse assignment is successfully deleted.
	time.Sleep(30 * time.Second)

	return nil
}

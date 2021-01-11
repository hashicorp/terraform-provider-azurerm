package lighthouse

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/lighthouse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/lighthouse/validate"
	resourceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	subscriptionValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLighthouseAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLighthouseAssignmentCreate,
		Read:   resourceLighthouseAssignmentRead,
		Delete: resourceLighthouseAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"lighthouse_definition_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LighthouseDefinitionID,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.Any(subscriptionValidate.SubscriptionID, resourceValidate.ResourceGroupID),
			},
		},
	}
}

func resourceLighthouseAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.AssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseAssignmentName := d.Get("name").(string)
	if lighthouseAssignmentName == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Lighthouse Assignment: %+v", err)
		}

		lighthouseAssignmentName = uuid
	}

	scope := d.Get("scope").(string)

	existing, err := client.Get(ctx, scope, lighthouseAssignmentName, utils.Bool(false))
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing Lighthouse Assignment %q (Scope %q): %+v", lighthouseAssignmentName, scope, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_lighthouse_assignment", *existing.ID)
	}

	parameters := managedservices.RegistrationAssignment{
		Properties: &managedservices.RegistrationAssignmentProperties{
			RegistrationDefinitionID: utils.String(d.Get("lighthouse_definition_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, scope, lighthouseAssignmentName, parameters); err != nil {
		return fmt.Errorf("creating Lighthouse Assignment %q (Scope %q): %+v", lighthouseAssignmentName, scope, err)
	}

	read, err := client.Get(ctx, scope, lighthouseAssignmentName, utils.Bool(false))
	if err != nil {
		return fmt.Errorf("retrieving Lighthouse Assessment %q (Scope %q): %+v", lighthouseAssignmentName, scope, err)
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("ID was nil or empty for Lighthouse Assignment %q ID (scope %q) ID", lighthouseAssignmentName, scope)
	}

	d.SetId(*read.ID)

	return resourceLighthouseAssignmentRead(d, meta)
}

func resourceLighthouseAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.AssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LighthouseAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Scope, id.Name, utils.Bool(false))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Lighthouse Assignment %q was not found (Scope %q)", id.Name, id.Scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Lighthouse Assignment %q (Scope %q): %+v", id.Name, id.Scope, err)
	}

	d.Set("name", resp.Name)
	d.Set("scope", id.Scope)

	if props := resp.Properties; props != nil {
		d.Set("lighthouse_definition_id", props.RegistrationDefinitionID)
	}

	return nil
}

func resourceLighthouseAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.AssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LighthouseAssignmentID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.Scope, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Lighthouse Assignment %q at Scope %q: %+v", id.Name, id.Scope, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"Deleted"},
		Refresh:    lighthouseAssignmentDeleteRefreshFunc(ctx, client, id.Scope, id.Name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Lighthouse Assignment %q (Scope %q) to be deleted: %s", id.Name, id.Scope, err)
	}

	return nil
}

func lighthouseAssignmentDeleteRefreshFunc(ctx context.Context, client *managedservices.RegistrationAssignmentsClient, scope string, lighthouseAssignmentName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		expandLighthouseDefinition := true
		res, err := client.Get(ctx, scope, lighthouseAssignmentName, &expandLighthouseDefinition)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("Error issuing read request in lighthouseAssignmentDeleteRefreshFunc to Lighthouse Assignment %q (Scope %q): %s", lighthouseAssignmentName, scope, err)
		}

		return res, "Deleting", nil
	}
}

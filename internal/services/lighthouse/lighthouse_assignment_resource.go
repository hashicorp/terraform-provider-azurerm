package lighthouse

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/lighthouse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/lighthouse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLighthouseAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLighthouseAssignmentCreate,
		Read:   resourceLighthouseAssignmentRead,
		Delete: resourceLighthouseAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LighthouseAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"lighthouse_definition_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LighthouseDefinitionID,
			},

			"scope": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.Any(commonids.ValidateSubscriptionID, commonids.ValidateResourceGroupID),
			},
		},
	}
}

func resourceLighthouseAssignmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.AssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseAssignmentName := d.Get("name").(string)
	if lighthouseAssignmentName == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Lighthouse Assignment: %+v", err)
		}

		lighthouseAssignmentName = uuid
	}

	id := parse.NewLighthouseAssignmentID(d.Get("scope").(string), lighthouseAssignmentName)
	existing, err := client.Get(ctx, id.Scope, id.Name, utils.Bool(false))
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	parameters := managedservices.RegistrationAssignment{
		Properties: &managedservices.RegistrationAssignmentProperties{
			RegistrationDefinitionID: utils.String(d.Get("lighthouse_definition_id").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.Scope, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLighthouseAssignmentRead(d, meta)
}

func resourceLighthouseAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("scope", id.Scope)

	if props := resp.Properties; props != nil {
		d.Set("lighthouse_definition_id", props.RegistrationDefinitionID)
	}

	return nil
}

func resourceLighthouseAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.AssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LighthouseAssignmentID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.Scope, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Lighthouse Assignment %q at Scope %q: %+v", id.Name, id.Scope, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"Deleted"},
		Refresh:    lighthouseAssignmentDeleteRefreshFunc(ctx, client, id.Scope, id.Name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Lighthouse Assignment %q (Scope %q) to be deleted: %s", id.Name, id.Scope, err)
	}

	return nil
}

func lighthouseAssignmentDeleteRefreshFunc(ctx context.Context, client *managedservices.RegistrationAssignmentsClient, scope string, lighthouseAssignmentName string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		expandLighthouseDefinition := true
		res, err := client.Get(ctx, scope, lighthouseAssignmentName, &expandLighthouseDefinition)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("issuing read request in lighthouseAssignmentDeleteRefreshFunc to Lighthouse Assignment %q (Scope %q): %s", lighthouseAssignmentName, scope, err)
		}

		return res, "Deleting", nil
	}
}

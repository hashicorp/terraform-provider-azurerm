package healthcare

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	workspace "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisWorkspaceCreateUpdate,
		Read:   resourceHealthcareApisWorkspaceRead,
		Update: resourceHealthcareApisWorkspaceCreateUpdate,
		Delete: resourceHealthcareApisWorkspaceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// todo: add the validation function
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareApisWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Workspace creation.")

	id := workspace.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_healthcareapis_workspace", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	parameters := workspace.Workspace{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisWorkspaceRead(d, meta)
}
func resourceHealthcareApisWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspace.ParseWorkspaceIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}
func resourceHealthcareApisWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspace.ParseWorkspaceIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(future.HttpResponse){
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return waitForHealthcareApiWorkspaceToBeDeleted(ctx, client, *id)
}

func waitForHealthcareApiWorkspaceToBeDeleted(ctx context.Context, client *workspace.WorkspacesClient, id workspace.WorkspaceId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: healthcareapiWorkspaceStateCodeRefreshFunc(ctx, client, id),
		Timeout: time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func healthcareapiWorkspaceStateCodeRefreshFunc(ctx context.Context, client *workspace.WorkspacesClient, id workspace.WorkspaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if res.HttpResponse != nil {
			log.Printf("Retrieving %s returned Status %d", id, res.HttpResponse.StatusCode)
		}

		if err != nil {
			if response.WasNotFound(res.HttpResponse){
				return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}

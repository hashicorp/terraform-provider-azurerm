package databricks

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-04-01-preview/accessconnector"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksAccessConnector() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatabricksAccessConnectorCreateUpdate,
		Read:   resourceDatabricksAccessConnectorRead,
		Update: resourceDatabricksAccessConnectorCreateUpdate,
		Delete: resourceDatabricksAccessConnectorDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := accessconnector.ParseAccessConnectorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccessConnectorName,
			},

			"location":            commonschema.Location(),
			"resource_group_name": commonschema.ResourceGroupName(),
			"identity":            commonschema.SystemAssignedIdentityRequired(),
			"tags":                commonschema.Tags(),
		},
	}
}

func resourceDatabricksAccessConnectorCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.AccessConnectorClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := accessconnector.NewAccessConnectorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_databricks_access_connector", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	identity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	accessConnector := accessconnector.AccessConnector{
		Location: location,
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, accessConnector); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDatabricksAccessConnectorRead(d, meta)
}

func resourceDatabricksAccessConnectorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.AccessConnectorClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accessconnector.ParseAccessConnectorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceDatabricksAccessConnectorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.AccessConnectorClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accessconnector.ParseAccessConnectorID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

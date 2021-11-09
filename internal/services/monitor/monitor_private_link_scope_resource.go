package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-07-01-preview/insights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorPrivateLinkScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorPrivateLinkScopeCreateUpdate,
		Read:   resourceMonitorPrivateLinkScopeRead,
		Update: resourceMonitorPrivateLinkScopeCreateUpdate,
		Delete: resourceMonitorPrivateLinkScopeDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PrivateLinkScopeID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PrivateLinkScopeName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorPrivateLinkScopeCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewPrivateLinkScopeID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_monitor_private_link_scope", id.ID())
		}
	}

	parameters := insights.AzureMonitorPrivateLinkScope{
		Location: utils.String("Global"),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorPrivateLinkScopeRead(d, meta)
}

func resourceMonitorPrivateLinkScopeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateLinkScopeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMonitorPrivateLinkScopeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateLinkScopeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

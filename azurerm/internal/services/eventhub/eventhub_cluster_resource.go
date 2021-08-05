package eventhub

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/eventhubsclusters"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceEventHubCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubClusterCreateUpdate,
		Read:   resourceEventHubClusterRead,
		Update: resourceEventHubClusterCreateUpdate,
		Delete: resourceEventHubClusterDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventhubsclusters.ClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			// You can't delete a cluster until at least 4 hours have passed from the initial creation.
			Delete: pluginsdk.DefaultTimeout(300 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^Dedicated_[1-9][0-9]*$`),
					"SKU name must match /^Dedicated_[1-9][0-9]*$/.",
				),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceEventHubClusterCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM EventHub Cluster creation.")

	id := eventhubsclusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.ClustersGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventhub_cluster", id.ID())
		}
	}

	cluster := eventhubsclusters.Cluster{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:     expandTags(d.Get("tags").(map[string]interface{})),
		Sku:      expandEventHubClusterSkuName(d.Get("sku_name").(string)),
	}

	if err := client.ClustersCreateOrUpdateThenPoll(ctx, id, cluster); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceEventHubClusterRead(d, meta)
}

func resourceEventHubClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventhubsclusters.ClusterID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.ClustersGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("sku_name", flattenEventHubClusterSkuName(model.Sku))
		d.Set("location", location.NormalizeNilable(model.Location))

		return tags.FlattenAndSet(d, flattenTags(model.Tags))
	}

	return nil
}

func resourceEventHubClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := eventhubsclusters.ClusterID(d.Id())
	if err != nil {
		return err
	}

	// The EventHub Cluster can't be deleted until four hours after creation so we'll keep retrying until it can be deleted.
	return pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutDelete), func() *pluginsdk.RetryError {
		future, err := client.ClustersDelete(ctx, *id)
		if err != nil {
			if response.WasNotFound(future.HttpResponse) {
				return nil
			}
			if strings.Contains(err.Error(), "Cluster cannot be deleted until four hours after its creation time") || future.HttpResponse.StatusCode == 429 {
				return pluginsdk.RetryableError(fmt.Errorf("expected eventhub cluster to be deleted but was in pending creation state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("deleting %s: %+v", *id, err))
		}

		if err := future.Poller.PollUntilDone(); err != nil {
			if response.WasNotFound(future.Poller.HttpResponse) {
				return nil
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("deleting %s: %+v", *id, err))
		}

		return nil
	}) //lintignore:R006
}

func expandEventHubClusterSkuName(skuName string) *eventhubsclusters.ClusterSku {
	if len(skuName) == 0 {
		return nil
	}

	name, capacity, err := azure.SplitSku(skuName)
	if err != nil {
		return nil
	}

	return &eventhubsclusters.ClusterSku{
		Name:     eventhubsclusters.ClusterSkuName(name),
		Capacity: utils.Int64(int64(capacity)),
	}
}

func flattenEventHubClusterSkuName(input *eventhubsclusters.ClusterSku) string {
	if input == nil || input.Capacity == nil {
		return ""
	}

	return fmt.Sprintf("%s_%d", string(input.Name), *input.Capacity)
}

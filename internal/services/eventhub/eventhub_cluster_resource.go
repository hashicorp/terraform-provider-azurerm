// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubsclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceEventHubCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubClusterCreateUpdate,
		Read:   resourceEventHubClusterRead,
		Update: resourceEventHubClusterCreateUpdate,
		Delete: resourceEventHubClusterDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventhubsclusters.ParseClusterID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,48}[a-zA-Z0-9])?$"),
					"The event hub name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 50 characters, and it must begin and end with a letter or number.",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^Dedicated_[1-9][0-9]*$`),
					"SKU name must match /^Dedicated_[1-9][0-9]*$/.",
				),
			},

			"tags": commonschema.Tags(),
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

	sku := expandEventHubClusterSkuName(d.Get("sku_name").(string))
	cluster := eventhubsclusters.Cluster{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku:      &sku,
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

	id, err := eventhubsclusters.ParseClusterID(d.Id())
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

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("sku_name", flattenEventHubClusterSkuName(model.Sku))
		d.Set("location", location.NormalizeNilable(model.Location))

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceEventHubClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := eventhubsclusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	// The EventHub Cluster can't be deleted until four hours after creation so we'll keep retrying until it can be deleted.
	return pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutDelete), func() *pluginsdk.RetryError {
		future, err := client.ClustersDelete(ctx, *id)
		if err != nil {
			if strings.Contains(err.Error(), "ClusterMoratoriumInEffect") || response.WasBadRequest(future.HttpResponse) || response.WasStatusCode(future.HttpResponse, http.StatusTooManyRequests) {
				return pluginsdk.RetryableError(fmt.Errorf("expected eventhub cluster to be deleted but was in pending creation state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("deleting %s: %+v", *id, err))
		}

		if err := future.Poller.PollUntilDone(ctx); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("deleting %s: %+v", *id, err))
		}

		return nil
	})
}

func expandEventHubClusterSkuName(skuName string) eventhubsclusters.ClusterSku {
	// "sku_name" is validated to be in this format above, and is required
	skuParts := strings.Split(skuName, "_")
	name := skuParts[0]
	capacity, _ := strconv.Atoi(skuParts[1])
	return eventhubsclusters.ClusterSku{
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

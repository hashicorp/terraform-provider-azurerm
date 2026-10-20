// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubClusterCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := eventhubsclusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	}

	sku := expandEventHubClusterSkuName(d.Get("sku_name").(string))
	cluster := eventhubsclusters.Cluster{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Sku:      &sku,
	}

	if d.IsNewResource() {
		if err := client.ClustersCreateOrUpdateCallbackThenPoll(ctx, id, cluster, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.ClustersCreateOrUpdateThenPoll(ctx, id, cluster); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceEventHubClusterRead(d, meta)
}

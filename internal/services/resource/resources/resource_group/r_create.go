// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceResourceGroupCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.ResourceGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewResourceGroupID(meta.(*clients.Client).Account.SubscriptionId, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(resourceGroupResourceName, id.ID())
		}
	}

	parameters := resourcegroups.ResourceGroup{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v := d.Get("managed_by").(string); v != "" {
		parameters.ManagedBy = pointer.To(v)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	// custom poller to account for replication delays in the eventual consistency responses of newly created RG resources
	pollerType := custompollers.NewResourceGroupCreatePoller(client, id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceResourceGroupRead(d, meta)
}

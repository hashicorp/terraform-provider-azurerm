// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_private_link_scoped_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2019-10-17-preview/privatelinkscopedresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorPrivateLinkScopedServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resourceGroupName := d.Get("resource_group_name").(string)
	scopeName := d.Get("scope_name").(string)
	if scopeResourceGroupName, ok := d.GetOk("scope_resource_group_name"); ok {
		resourceGroupName = scopeResourceGroupName.(string)
	}
	if scopeSubscriptionId, ok := d.GetOk("scope_subscription_id"); ok {
		subscriptionId = scopeSubscriptionId.(string)
	}
	client := meta.(*clients.Client).Monitor.PrivateLinkScopedResourcesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if scopeIdRaw, ok := d.GetOk("scope_resource_id"); ok {
		scopeId, err := privatelinkscopesapis.ParsePrivateLinkScopeID(scopeIdRaw.(string))
		if err != nil {
			return fmt.Errorf("parsing scope ID: %+v", err)
		}

		subscriptionId = scopeId.SubscriptionId
		resourceGroupName = scopeId.ResourceGroupName
		scopeName = scopeId.PrivateLinkScopeName
	}

	id := privatelinkscopedresources.NewScopedResourceID(subscriptionId, resourceGroupName, scopeName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_private_link_scoped_service", id.ID())
		}
	}

	parameters := privatelinkscopedresources.ScopedResource{
		Properties: &privatelinkscopedresources.ScopedResourceProperties{
			LinkedResourceId: pointer.To(d.Get("linked_resource_id").(string)),
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceMonitorPrivateLinkScopedServiceRead(d, meta)
}

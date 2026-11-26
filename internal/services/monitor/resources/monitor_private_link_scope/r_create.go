// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_private_link_scope

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorPrivateLinkScopeCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := privatelinkscopesapis.NewPrivateLinkScopeID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.PrivateLinkScopesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_monitor_private_link_scope", id.ID())
			}
		}
	}

	ingestionAccessMode := privatelinkscopesapis.AccessMode(d.Get("ingestion_access_mode").(string))
	queryaccessMode := privatelinkscopesapis.AccessMode(d.Get("query_access_mode").(string))

	parameters := privatelinkscopesapis.AzureMonitorPrivateLinkScope{
		Name:     &name,
		Location: "Global",
		Tags:     pluginsdk.ExpandPtrMapStringString(d.Get("tags").(map[string]any)),
		Properties: privatelinkscopesapis.AzureMonitorPrivateLinkScopeProperties{
			AccessModeSettings: privatelinkscopesapis.AccessModeSettings{
				IngestionAccessMode: ingestionAccessMode,
				QueryAccessMode:     queryaccessMode,
			},
		},
	}

	if _, err := client.PrivateLinkScopesCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorPrivateLinkScopeRead(d, meta)
}

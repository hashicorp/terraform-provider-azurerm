// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_private_link_scoped_service

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2019-10-17-preview/privatelinkscopedresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorPrivateLinkScopedServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopedResourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkscopedresources.ParseScopedResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ScopedResourceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("scope_name", id.PrivateLinkScopeName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("linked_resource_id", normalizeLinkedResourceId(props.LinkedResourceId))
		}
	}

	return nil
}

func normalizeLinkedResourceId(input *string) *string {
	if input == nil {
		return input
	}

	if resourceId, err := componentsapis.ParseComponentIDInsensitively(*input); err == nil {
		normalizedId := resourceId.ID()
		return &normalizedId
	}
	if resourceId, err := workspaces.ParseWorkspaceIDInsensitively(*input); err == nil {
		normalizedId := resourceId.ID()
		return &normalizedId
	}
	if resourceId, err := datacollectionendpoints.ParseDataCollectionEndpointIDInsensitively(*input); err == nil {
		normalizedId := resourceId.ID()
		return &normalizedId
	}

	return input
}

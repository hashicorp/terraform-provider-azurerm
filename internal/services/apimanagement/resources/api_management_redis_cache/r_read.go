// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_redis_cache

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/cache"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementRedisCacheRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.CacheClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceManagerEndpoint, ok := meta.(*clients.Client).Account.Environment.ResourceManager.Endpoint()
	if !ok {
		return fmt.Errorf("could not determine Resource Manager endpoint suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	id, err := cache.ParseCacheID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] apimanagement %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("name", id.CacheId)
	d.Set("api_management_id", apimanagementservice.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName).ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))

			cacheId := ""
			if props.ResourceId != nil {
				// correct the resourceID issue: "https://management.azure.com//subscriptions/xx/resourceGroups/xx/providers/Microsoft.Cache/Redis/xx"
				cacheId = strings.TrimPrefix(*props.ResourceId, *resourceManagerEndpoint)
			}
			d.Set("redis_cache_id", cacheId)
			d.Set("cache_location", props.UseFromLocation)
		}
	}
	return nil
}

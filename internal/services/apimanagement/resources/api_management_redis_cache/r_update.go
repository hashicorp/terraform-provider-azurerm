// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_redis_cache

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/cache"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementRedisCacheCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.CacheClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceManagerEndpoint, ok := meta.(*clients.Client).Account.Environment.ResourceManager.Endpoint()
	if !ok {
		return fmt.Errorf("could not determine Resource Manager endpoint suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	name := d.Get("name").(string)
	apimId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return err
	}
	id := cache.NewCacheID(subscriptionId, apimId.ResourceGroupName, apimId.ServiceName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_redis_cache", id.ID())
			}
		}
	}

	parameters := cache.CacheContract{
		Properties: &cache.CacheContractProperties{
			ConnectionString: d.Get("connection_string").(string),
			UseFromLocation:  location.Normalize(d.Get("cache_location").(string)),
		},
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		parameters.Properties.Description = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("redis_cache_id"); ok && v.(string) != "" {
		parameters.Properties.ResourceId = pointer.To(*resourceManagerEndpoint + v.(string))
	}

	// here we use "PUT" for updating, because `description` is not allowed to be empty string, Then we could not update to remove `description` by `PATCH`
	if _, err := client.CreateOrUpdate(ctx, id, parameters, cache.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/ updating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementRedisCacheRead(d, meta)
}

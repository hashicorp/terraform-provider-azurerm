// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/cache"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementRedisCache() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementRedisCacheCreateUpdate,
		Read:   resourceApiManagementRedisCacheRead,
		Update: resourceApiManagementRedisCacheCreateUpdate,
		Delete: resourceApiManagementRedisCacheDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cache.ParseCacheID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementChildName,
			},

			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: apimanagementservice.ValidateServiceID,
			},

			"connection_string": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"redis_cache_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"cache_location": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Default:          "default",
				ValidateFunc:     validate.RedisCacheLocation,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},
		},
	}
}

func resourceApiManagementRedisCacheCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementRedisCacheRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementRedisCacheDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.CacheClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cache.ParseCacheID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id, cache.DeleteOperationOptions{IfMatch: pointer.To("*")}); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	return nil
}

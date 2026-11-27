// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_redis_cache

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/cache"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementRedisCacheDelete(d *pluginsdk.ResourceData, meta any) error {
	client := apimanagementClient.From(meta).CacheClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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

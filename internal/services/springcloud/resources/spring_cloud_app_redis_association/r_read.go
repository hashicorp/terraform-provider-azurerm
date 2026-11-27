// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_redis_association

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudAppRedisAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).BindingsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform.ParseBindingIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.BindingName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Spring Cloud App Association %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.Set("name", id.BindingName)
	d.Set("spring_cloud_app_id", appplatform.NewAppID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.AppName).ID())
	if props := resp.Properties; props != nil {
		d.Set("redis_cache_id", props.ResourceID)

		enableSSL := "false"
		if v, ok := props.BindingParameters[springCloudAppRedisAssociationKeySSL]; ok && v != nil {
			enableSSL = *v
		}
		d.Set("ssl_enabled", strings.EqualFold(enableSSL, "true"))
	}
	return nil
}

const springCloudAppRedisAssociationKeySSL = "useSsl"

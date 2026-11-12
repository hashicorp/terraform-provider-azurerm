// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisCacheRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Redis.RedisResourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redisresources.ParseRediID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.RedisGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceRedisCacheFlatten(ctx, meta.(*clients.Client), d, id, resp.Model, true)
}

var skuWeight = map[string]int8{
	"Basic":    1,
	"Standard": 2,
	"Premium":  3,
}

func redisStateRefreshFunc(ctx context.Context, client *redisresources.RedisResourcesClient, id redisresources.RediId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.RedisGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for status of %s: %+v", id, err)
		}

		provisioningState := ""
		if model := res.Model; model != nil && model.Properties.ProvisioningState != nil {
			provisioningState = string(*res.Model.Properties.ProvisioningState)
		}
		if provisioningState == "" {
			return nil, "", fmt.Errorf("polling for status of %s: `provisioningState` was nil", id)
		}

		return res, provisioningState, nil
	}
}

func isAuthRequiredAsBool(notRequired string) bool {
	value := strings.ToLower(notRequired)
	output := map[string]bool{
		"yes": false,
		"no":  true,
	}
	return output[value]
}

func isAuthNotRequiredAsString(authRequired bool) string {
	output := map[bool]string{
		true:  "no",
		false: "yes",
	}
	return output[authRequired]
}

func getRedisConnectionString(redisHostName string, sslPort int64, accessKey string, enableSslPort bool) string {
	return fmt.Sprintf("%s:%d,password=%s,ssl=%t,abortConnect=False", redisHostName, sslPort, accessKey, enableSslPort)
}

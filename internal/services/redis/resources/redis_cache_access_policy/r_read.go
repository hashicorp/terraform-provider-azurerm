// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache_access_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := rediscacheaccesspolicies.ParseAccessPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Redis.CacheAccessPoliciesClient

			resp, err := client.AccessPolicyGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving Redis Cache Access Policy %s: %+v", *id, err)
			}

			state := RedisCacheAccessPolicyResourceModel{}

			if model := resp.Model; model != nil {
				if model.Name != nil {
					state.Name = *model.Name
				}
				state.RedisCacheID = rediscacheaccesspolicies.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName).ID()
				if model.Properties != nil {
					state.Permissions = model.Properties.Permissions
				}
			}

			return metadata.Encode(&state)
		},
	}
}

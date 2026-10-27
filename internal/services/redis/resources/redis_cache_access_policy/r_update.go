// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache_access_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Redis.CacheAccessPoliciesClient

			var state RedisCacheAccessPolicyResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := rediscacheaccesspolicies.ParseAccessPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.AccessPolicyGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			model := *existing.Model
			if metadata.ResourceData.HasChange("permissions") {
				model.Properties.Permissions = state.Permissions
			}

			locks.ByID(state.RedisCacheID)
			defer locks.UnlockByID(state.RedisCacheID)

			if err := client.AccessPolicyCreateUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

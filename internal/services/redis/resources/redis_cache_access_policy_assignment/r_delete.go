// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache_access_policy_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RedisCacheAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.Redis.CacheAccessPolicyAssignmentsClient
			id, err := rediscacheaccesspolicyassignments.ParseAccessPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(model.RedisCacheID)
			defer locks.UnlockByID(model.RedisCacheID)

			if err := client.AccessPolicyAssignmentDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting Redis Cache Access Policy Assignment %s in Redis Cache %s in resource group %s: %s", id.AccessPolicyAssignmentName, id.RedisName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}

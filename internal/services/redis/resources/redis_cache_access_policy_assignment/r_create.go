// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache_access_policy_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RedisCacheAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.Redis.CacheAccessPolicyAssignmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			redisId, err := rediscacheaccesspolicyassignments.ParseRediID(model.RedisCacheID)
			if err != nil {
				return err
			}
			id := rediscacheaccesspolicyassignments.NewAccessPolicyAssignmentID(subscriptionId, redisId.ResourceGroupName, redisId.RedisName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.AccessPolicyAssignmentGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			createInput := rediscacheaccesspolicyassignments.RedisCacheAccessPolicyAssignment{
				Name: &model.Name,
				Properties: &rediscacheaccesspolicyassignments.RedisCacheAccessPolicyAssignmentProperties{
					AccessPolicyName: model.AccessPolicyName,
					ObjectId:         model.ObjectID,
					ObjectIdAlias:    model.ObjectIDAlias,
				},
			}

			locks.ByID(model.RedisCacheID)
			defer locks.UnlockByID(model.RedisCacheID)

			if err := client.AccessPolicyAssignmentCreateUpdateCallbackThenPoll(ctx, id, createInput, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/redis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RedisCacheAccessPolicyAssignmentResource struct {
}

var _ sdk.Resource = RedisCacheAccessPolicyAssignmentResource{}

type RedisCacheAccessPolicyAssignmentResourceModel struct {
	Name             string `tfschema:"name"`
	RedisCacheID     string `tfschema:"redis_cache_id"`
	AccessPolicyName string `tfschema:"access_policy_name"`
	ObjectID         string `tfschema:"object_id"`
	ObjectIDAlias    string `tfschema:"object_id_alias"`
}

func (r RedisCacheAccessPolicyAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"redis_cache_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: redis.ValidateRediID,
		},
		"access_policy_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"object_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"object_id_alias": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r RedisCacheAccessPolicyAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r RedisCacheAccessPolicyAssignmentResource) ModelObject() interface{} {
	return &RedisCacheAccessPolicyAssignmentResourceModel{}
}

func (r RedisCacheAccessPolicyAssignmentResource) ResourceType() string {
	return "azurerm_redis_cache_access_policy_assignment"
}

func (r RedisCacheAccessPolicyAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return redis.ValidateAccessPolicyAssignmentID
}

func (r RedisCacheAccessPolicyAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RedisCacheAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.Redis.Redis
			subscriptionId := metadata.Client.Account.SubscriptionId

			redisId, err := redis.ParseRediID(model.RedisCacheID)
			if err != nil {
				return err
			}
			id := redis.NewAccessPolicyAssignmentID(subscriptionId, redisId.ResourceGroupName, redisId.RedisName, model.Name)

			existing, err := client.AccessPolicyAssignmentGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			createInput := redis.RedisCacheAccessPolicyAssignment{
				Name: &model.Name,
				Properties: &redis.RedisCacheAccessPolicyAssignmentProperties{
					AccessPolicyName: model.AccessPolicyName,
					ObjectId:         model.ObjectID,
					ObjectIdAlias:    model.ObjectIDAlias,
				},
			}

			locks.ByID(model.RedisCacheID)
			defer locks.UnlockByID(model.RedisCacheID)

			if err := client.AccessPolicyAssignmentCreateUpdateThenPoll(ctx, id, createInput); err != nil {
				return fmt.Errorf("failed to create Redis Cache Access Policy Assignment %s in Redis Cache %s in resource group %s: %s", model.Name, redisId.RedisName, redisId.ResourceGroupName, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RedisCacheAccessPolicyAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := redis.ParseAccessPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Redis.Redis

			resp, err := client.AccessPolicyAssignmentGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving Redis Cache Access Policy Assignment %s: %+v", *id, err)
			}

			state := RedisCacheAccessPolicyAssignmentResourceModel{}

			if model := resp.Model; model != nil {
				if model.Name != nil {
					state.Name = *model.Name
				}
				state.RedisCacheID = redis.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName).ID()
				if model.Properties != nil {
					state.AccessPolicyName = model.Properties.AccessPolicyName
					state.ObjectID = model.Properties.ObjectId
					state.ObjectIDAlias = model.Properties.ObjectIdAlias
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RedisCacheAccessPolicyAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RedisCacheAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.Redis.Redis
			id, err := redis.ParseAccessPolicyAssignmentID(metadata.ResourceData.Id())
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

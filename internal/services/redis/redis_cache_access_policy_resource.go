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
)

type RedisCacheAccessPolicyResource struct{}

var _ sdk.ResourceWithUpdate = RedisCacheAccessPolicyResource{}

type RedisCacheAccessPolicyResourceModel struct {
	Name         string `tfschema:"name"`
	RedisCacheID string `tfschema:"redis_cache_id"`
	Permissions  string `tfschema:"permissions"`
}

func (r RedisCacheAccessPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"redis_cache_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: redis.ValidateRediID,
		},
		"permissions": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (r RedisCacheAccessPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r RedisCacheAccessPolicyResource) ModelObject() interface{} {
	return &RedisCacheAccessPolicyResourceModel{}
}

func (r RedisCacheAccessPolicyResource) ResourceType() string {
	return "azurerm_redis_cache_access_policy"
}

func (r RedisCacheAccessPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return redis.ValidateAccessPolicyID
}

func (r RedisCacheAccessPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RedisCacheAccessPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.Redis.Redis
			subscriptionId := metadata.Client.Account.SubscriptionId

			redisId, err := redis.ParseRediID(model.RedisCacheID)
			if err != nil {
				return fmt.Errorf("parsing Redis Cache ID (%s): %+v", model.RedisCacheID, err)
			}
			id := redis.NewAccessPolicyID(subscriptionId, redisId.ResourceGroupName, redisId.RedisName, model.Name)

			existing, err := client.AccessPolicyGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			policyTypeCustom := redis.AccessPolicyTypeCustom

			createInput := redis.RedisCacheAccessPolicy{
				Name: &model.Name,
				Properties: &redis.RedisCacheAccessPolicyProperties{
					Permissions: model.Permissions,
					Type:        &policyTypeCustom,
				},
			}

			locks.ByID(model.RedisCacheID)
			defer locks.UnlockByID(model.RedisCacheID)

			if err := client.AccessPolicyCreateUpdateThenPoll(ctx, id, createInput); err != nil {
				return fmt.Errorf("failed to create Redis Cache Access Policy %s in Redis Cache %s in resource group %s: %s", model.Name, redisId.RedisName, redisId.ResourceGroupName, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RedisCacheAccessPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Redis.Redis

			var state RedisCacheAccessPolicyResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := redis.ParseAccessPolicyID(metadata.ResourceData.Id())
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

func (r RedisCacheAccessPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := redis.ParseAccessPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Redis.Redis

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
				state.RedisCacheID = redis.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName).ID()
				if model.Properties != nil {
					state.Permissions = model.Properties.Permissions
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RedisCacheAccessPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RedisCacheAccessPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.Redis.Redis
			id, err := redis.ParseAccessPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			locks.ByID(model.RedisCacheID)
			defer locks.UnlockByID(model.RedisCacheID)

			if err := client.AccessPolicyDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting Redis Cache Access Policy %s in Redis Cache %s in resource group %s: %s", id.AccessPolicyName, id.RedisName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}

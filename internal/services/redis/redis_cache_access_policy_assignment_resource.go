package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-08-01/redis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RedisCacheAccessPolicyAssignmentResource struct {
}

var _ sdk.ResourceWithUpdate = RedisCacheAccessPolicyAssignmentResource{}

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
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"redis_cache_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"access_policy_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"object_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"object_id_alias": {
			Type:     pluginsdk.TypeString,
			Required: true,
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
				return fmt.Errorf("parsing Redis Cache ID (%s): %+v", model.RedisCacheID, err)
			}
			id := redis.NewAccessPolicyAssignmentID(subscriptionId, redisId.ResourceGroupName, redisId.RedisName, model.Name)

			createInput := redis.RedisCacheAccessPolicyAssignment{
				Name: &model.Name,
				Properties: &redis.RedisCacheAccessPolicyAssignmentProperties{
					AccessPolicyName: model.AccessPolicyName,
					ObjectId:         model.ObjectID,
					ObjectIdAlias:    model.ObjectIDAlias,
				},
			}
			if err := client.AccessPolicyAssignmentCreateUpdateThenPoll(ctx, id, createInput); err != nil {
				return fmt.Errorf("failed to create Redis Cache Access Policy Assignment %s in Redis Cache %s in resource group %s: %s", model.Name, redisId.RedisName, redisId.ResourceGroupName, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RedisCacheAccessPolicyAssignmentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Redis.Redis
			id, err := redis.ParseAccessPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state RedisCacheAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.AccessPolicyAssignmentGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Redis Cache Access Policy Assignment %s: %v", id, err)
			}

			updateInput := redis.RedisCacheAccessPolicyAssignment{
				Name: existing.Model.Name,
				Properties: &redis.RedisCacheAccessPolicyAssignmentProperties{
					AccessPolicyName: existing.Model.Properties.AccessPolicyName,
					ObjectId:         existing.Model.Properties.ObjectId,
					ObjectIdAlias:    existing.Model.Properties.ObjectIdAlias,
				},
			}

			if metadata.ResourceData.HasChange("name") {
				updateInput.Name = &state.Name
			}
			if metadata.ResourceData.HasChange("access_policy_name") {
				updateInput.Properties.AccessPolicyName = state.AccessPolicyName
			}
			if metadata.ResourceData.HasChange("object_id") {
				updateInput.Properties.ObjectId = state.ObjectID
			}
			if metadata.ResourceData.HasChange("object_id_alias") {
				updateInput.Properties.ObjectIdAlias = state.ObjectIDAlias
			}
			if err := client.AccessPolicyAssignmentCreateUpdateThenPoll(ctx, *id, updateInput); err != nil {
				return fmt.Errorf("failed to update Redis Cache Access Policy Assignment: %s: %+v", id.AccessPolicyAssignmentName, err)
			}

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
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RedisCacheAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.Redis.Redis
			id, err := redis.ParseAccessPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			if _, err := client.AccessPolicyAssignmentDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting Redis Cache Access Policy Assignment %s in Redis Cache %s in resource group %s: %s", id.AccessPolicyAssignmentName, id.RedisName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}

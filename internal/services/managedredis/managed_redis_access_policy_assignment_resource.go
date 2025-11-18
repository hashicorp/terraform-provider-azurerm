// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagedRedisAccessPolicyAssignmentResource struct{}

var _ sdk.Resource = ManagedRedisAccessPolicyAssignmentResource{}

type ManagedRedisAccessPolicyAssignmentResourceModel struct {
	ManagedRedisID string `tfschema:"managed_redis_id"`
	DatabaseName   string `tfschema:"database_name"`
	ObjectID       string `tfschema:"object_id"`
}

func (r ManagedRedisAccessPolicyAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_redis_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
		},

		"database_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  defaultDatabaseName,
		},

		"object_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r ManagedRedisAccessPolicyAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedRedisAccessPolicyAssignmentResource) ModelObject() interface{} {
	return &ManagedRedisAccessPolicyAssignmentResourceModel{}
}

func (r ManagedRedisAccessPolicyAssignmentResource) ResourceType() string {
	return "azurerm_managed_redis_access_policy_assignment"
}

func (r ManagedRedisAccessPolicyAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return databases.ValidateAccessPolicyAssignmentID
}

func (r ManagedRedisAccessPolicyAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedRedisAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(model.ManagedRedisID)
			if err != nil {
				return fmt.Errorf("parsing managed redis ID: %+v", err)
			}

			// Access policy assignments are created on the specified database
			// Use object_id as the assignment name to ensure one assignment per user per database
			dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, model.DatabaseName)
			id := databases.NewAccessPolicyAssignmentID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, model.DatabaseName, model.ObjectID)

			existing, err := client.AccessPolicyAssignmentGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// Check that the database exists
			dbResp, err := client.Get(ctx, dbId)
			if err != nil {
				if response.WasNotFound(dbResp.HttpResponse) {
					return fmt.Errorf("managed Redis database %s was not found", dbId)
				}
				return fmt.Errorf("retrieving %s: %+v", dbId, err)
			}

			createInput := databases.AccessPolicyAssignment{
				Name: pointer.To(model.ObjectID),
				Properties: &databases.AccessPolicyAssignmentProperties{
					AccessPolicyName: "default",
					User: databases.AccessPolicyAssignmentPropertiesUser{
						ObjectId: pointer.To(model.ObjectID),
					},
				},
			}

			locks.ByID(model.ManagedRedisID)
			defer locks.UnlockByID(model.ManagedRedisID)

			if err := client.AccessPolicyAssignmentCreateUpdateThenPoll(ctx, id, createInput); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagedRedisAccessPolicyAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := databases.ParseAccessPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.ManagedRedis.DatabaseClient

			resp, err := client.AccessPolicyAssignmentGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ManagedRedisAccessPolicyAssignmentResourceModel{}

			if model := resp.Model; model != nil {
				clusterId := redisenterprise.NewRedisEnterpriseID(id.SubscriptionId, id.ResourceGroupName, id.RedisEnterpriseName)
				state.ManagedRedisID = clusterId.ID()
				state.DatabaseName = id.DatabaseName

				if props := model.Properties; props != nil {
					if props.User.ObjectId != nil {
						state.ObjectID = *props.User.ObjectId
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedRedisAccessPolicyAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedRedisAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ManagedRedis.DatabaseClient
			id, err := databases.ParseAccessPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(model.ManagedRedisID)
			defer locks.UnlockByID(model.ManagedRedisID)

			if err := client.AccessPolicyAssignmentDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

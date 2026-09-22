// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis_access_policy_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedRedisAccessPolicyAssignmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(model.ManagedRedisID)
			if err != nil {
				return err
			}

			// Access policy assignments are created on the specified database
			// Use object_id as the assignment name to ensure one assignment per user per database
			dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)
			id := databases.NewAccessPolicyAssignmentID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName, model.ObjectID)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.AccessPolicyAssignmentGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

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

			if err := client.AccessPolicyAssignmentCreateUpdateCallbackThenPoll(ctx, id, createInput, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

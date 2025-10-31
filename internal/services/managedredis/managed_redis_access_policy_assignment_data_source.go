// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedRedisAccessPolicyAssignmentDataSource struct{}

var _ sdk.DataSource = ManagedRedisAccessPolicyAssignmentDataSource{}

type ManagedRedisAccessPolicyAssignmentDataSourceModel struct {
	Name             string `tfschema:"name"`
	ManagedRedisName string `tfschema:"managed_redis_name"`
	ResourceGroup    string `tfschema:"resource_group_name"`
	AccessPolicyName string `tfschema:"access_policy_name"`
	ObjectID         string `tfschema:"object_id"`
}

func (r ManagedRedisAccessPolicyAssignmentDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ManagedRedisAccessPolicyAssignmentName,
		},

		"managed_redis_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ManagedRedisClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r ManagedRedisAccessPolicyAssignmentDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"access_policy_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"object_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagedRedisAccessPolicyAssignmentDataSource) ModelObject() interface{} {
	return &ManagedRedisAccessPolicyAssignmentDataSourceModel{}
}

func (r ManagedRedisAccessPolicyAssignmentDataSource) ResourceType() string {
	return "azurerm_managed_redis_access_policy_assignment"
}

func (r ManagedRedisAccessPolicyAssignmentDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state ManagedRedisAccessPolicyAssignmentDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ManagedRedis.DatabaseClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := databases.NewAccessPolicyAssignmentID(subscriptionId, state.ResourceGroup, state.ManagedRedisName, defaultDatabaseName, state.Name)

			resp, err := client.AccessPolicyAssignmentGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.AccessPolicyName = props.AccessPolicyName
					if props.User.ObjectId != nil {
						state.ObjectID = *props.User.ObjectId
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_role_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.MongoRBACClient

			id, err := mongorbacs.ParseMongodbRoleDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.MongoDBResourcesGetMongoRoleDefinition(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := CosmosDbMongoRoleDefinitionResourceModel{}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					databaseId := cosmosdb.NewMongodbDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, pointer.From(properties.DatabaseName))

					state.CosmosMongoDatabaseId = databaseId.ID()
					state.RoleName = pointer.From(properties.RoleName)
					state.Privileges = flattenPrivilege(properties.Privileges)
					state.InheritedRoleNames = flattenInheritedRoles(properties.Roles)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

type Privilege struct {
	Actions  []string        `tfschema:"actions"`
	Resource []ResourceModel `tfschema:"resource"`
}

type ResourceModel struct {
	CollectionName string `tfschema:"collection_name"`
	DbName         string `tfschema:"db_name"`
}

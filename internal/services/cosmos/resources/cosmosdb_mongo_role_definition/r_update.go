// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_role_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.MongoRBACClient

			id, err := mongorbacs.ParseMongodbRoleDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
			defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

			var model CosmosDbMongoRoleDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.MongoDBResourcesGetMongoRoleDefinition(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			databaseId, err := cosmosdb.ParseMongodbDatabaseID(model.CosmosMongoDatabaseId)
			if err != nil {
				return err
			}

			roleType := mongorbacs.MongoRoleDefinitionTypeCustomRole
			parameters := mongorbacs.MongoRoleDefinitionCreateUpdateParameters{
				Properties: &mongorbacs.MongoRoleDefinitionResource{
					DatabaseName: pointer.To(databaseId.MongodbDatabaseName),
					RoleName:     pointer.To(model.RoleName),
					Roles:        expandInheritedRoles(model.InheritedRoleNames, databaseId.MongodbDatabaseName),
					Privileges:   expandPrivilege(model.Privileges),
					Type:         pointer.To(roleType),
				},
			}

			if err := client.MongoDBResourcesCreateUpdateMongoRoleDefinitionThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

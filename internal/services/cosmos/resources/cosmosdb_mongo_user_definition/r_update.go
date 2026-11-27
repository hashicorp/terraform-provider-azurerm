// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_user_definition

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

			id, err := mongorbacs.ParseMongodbUserDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
			defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

			var model CosmosDbMongoUserDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.MongoDBResourcesGetMongoUserDefinition(ctx, *id)
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

			parameters := mongorbacs.MongoUserDefinitionCreateUpdateParameters{
				Properties: properties.Properties,
			}

			if metadata.ResourceData.HasChange("password") {
				parameters.Properties.Password = pointer.To(model.Password)
			}

			if metadata.ResourceData.HasChange("inherited_role_names") {
				parameters.Properties.Roles = expandInheritedRole(model.InheritedRoleNames, databaseId.MongodbDatabaseName)
			}

			if err := client.MongoDBResourcesCreateUpdateMongoUserDefinitionThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

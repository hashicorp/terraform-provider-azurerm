// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_user_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbMongoUserDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.MongoRBACClient
			databaseId, err := cosmosdb.ParseMongodbDatabaseID(model.CosmosMongoDatabaseId)
			if err != nil {
				return err
			}

			mongoUserDefinitionId := fmt.Sprintf("%s.%s", databaseId.MongodbDatabaseName, model.Username)
			id := mongorbacs.NewMongodbUserDefinitionID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.DatabaseAccountName, mongoUserDefinitionId)

			locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
			defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.MongoDBResourcesGetMongoUserDefinition(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := mongorbacs.MongoUserDefinitionCreateUpdateParameters{
				Properties: &mongorbacs.MongoUserDefinitionResource{
					DatabaseName: pointer.To(databaseId.MongodbDatabaseName),
					Mechanisms:   pointer.To("SCRAM-SHA-256"),
					Password:     pointer.To(model.Password),
					UserName:     pointer.To(model.Username),
					Roles:        expandInheritedRole(model.InheritedRoleNames, databaseId.MongodbDatabaseName),
				},
			}

			if err := client.MongoDBResourcesCreateUpdateMongoUserDefinitionCallbackThenPoll(ctx, id, properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandInheritedRole(input []string, dbName string) *[]mongorbacs.Role {
	if len(input) == 0 {
		return nil
	}

	result := make([]mongorbacs.Role, 0)

	for _, v := range input {
		role := mongorbacs.Role{
			Db:   pointer.To(dbName),
			Role: pointer.To(v),
		}

		result = append(result, role)
	}

	return &result
}

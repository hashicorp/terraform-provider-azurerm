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
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbMongoRoleDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.MongoRBACClient
			databaseId, err := cosmosdb.ParseMongodbDatabaseID(model.CosmosMongoDatabaseId)
			if err != nil {
				return err
			}

			mongoRoleDefinitionId := fmt.Sprintf("%s.%s", databaseId.MongodbDatabaseName, model.RoleName)
			id := mongorbacs.NewMongodbRoleDefinitionID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.DatabaseAccountName, mongoRoleDefinitionId)

			locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
			defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.MongoDBResourcesGetMongoRoleDefinition(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			roleType := mongorbacs.MongoRoleDefinitionTypeCustomRole
			parameters := mongorbacs.MongoRoleDefinitionCreateUpdateParameters{
				Properties: &mongorbacs.MongoRoleDefinitionResource{
					DatabaseName: pointer.To(databaseId.MongodbDatabaseName),
					Roles:        expandInheritedRoles(model.InheritedRoleNames, databaseId.MongodbDatabaseName),
					RoleName:     pointer.To(model.RoleName),
					Privileges:   expandPrivilege(model.Privileges),
					Type:         pointer.To(roleType),
				},
			}

			if err := client.MongoDBResourcesCreateUpdateMongoRoleDefinitionCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandPrivilege(input []Privilege) *[]mongorbacs.Privilege {
	if len(input) == 0 {
		return nil
	}

	result := make([]mongorbacs.Privilege, 0, len(input))
	for _, v := range input {
		output := mongorbacs.Privilege{
			Actions:  pointer.To(v.Actions),
			Resource: expandResource(v.Resource),
		}

		result = append(result, output)
	}

	return &result
}

func expandResource(input []ResourceModel) *mongorbacs.PrivilegeResource {
	if len(input) == 0 {
		return nil
	}

	privilegeResource := &input[0]
	result := mongorbacs.PrivilegeResource{}

	if privilegeResource.CollectionName != "" {
		result.Collection = pointer.To(privilegeResource.CollectionName)
	}

	if privilegeResource.DbName != "" {
		result.Db = pointer.To(privilegeResource.DbName)
	}

	return &result
}

func expandInheritedRoles(input []string, dbName string) *[]mongorbacs.Role {
	if len(input) == 0 {
		return nil
	}

	result := make([]mongorbacs.Role, 0, len(input))
	for _, v := range input {
		inheritedRole := mongorbacs.Role{
			Db:   pointer.To(dbName),
			Role: pointer.To(v),
		}

		result = append(result, inheritedRole)
	}

	return &result
}

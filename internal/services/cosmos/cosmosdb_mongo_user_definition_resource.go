// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-11-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CosmosDbMongoUserDefinitionResourceModel struct {
	CosmosMongoDatabaseId string   `tfschema:"cosmos_mongo_database_id"`
	Username              string   `tfschema:"username"`
	Password              string   `tfschema:"password"`
	InheritedRoleNames    []string `tfschema:"inherited_role_names"`
}

type CosmosDbMongoUserDefinitionResource struct{}

var _ sdk.ResourceWithUpdate = CosmosDbMongoUserDefinitionResource{}

func (r CosmosDbMongoUserDefinitionResource) ResourceType() string {
	return "azurerm_cosmosdb_mongo_user_definition"
}

func (r CosmosDbMongoUserDefinitionResource) ModelObject() interface{} {
	return &CosmosDbMongoUserDefinitionResourceModel{}
}

func (r CosmosDbMongoUserDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return mongorbacs.ValidateMongodbUserDefinitionID
}

func (r CosmosDbMongoUserDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cosmos_mongo_database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.MongodbDatabaseID,
		},

		"username": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"inherited_role_names": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r CosmosDbMongoUserDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CosmosDbMongoUserDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbMongoUserDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.MongoRBACClient
			databaseId, err := parse.MongodbDatabaseID(model.CosmosMongoDatabaseId)
			if err != nil {
				return err
			}

			mongoUserDefinitionId := fmt.Sprintf("%s.%s", databaseId.Name, model.Username)
			id := mongorbacs.NewMongodbUserDefinitionID(databaseId.SubscriptionId, databaseId.ResourceGroup, databaseId.DatabaseAccountName, mongoUserDefinitionId)

			locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
			defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

			existing, err := client.MongoDBResourcesGetMongoUserDefinition(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := mongorbacs.MongoUserDefinitionCreateUpdateParameters{
				Properties: &mongorbacs.MongoUserDefinitionResource{
					DatabaseName: pointer.To(databaseId.Name),
					Mechanisms:   pointer.To("SCRAM-SHA-256"),
					Password:     pointer.To(model.Password),
					UserName:     pointer.To(model.Username),
					Roles:        expandInheritedRole(model.InheritedRoleNames, databaseId.Name),
				},
			}

			if err := client.MongoDBResourcesCreateUpdateMongoUserDefinitionThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CosmosDbMongoUserDefinitionResource) Update() sdk.ResourceFunc {
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

			databaseId, err := parse.MongodbDatabaseID(model.CosmosMongoDatabaseId)
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
				parameters.Properties.Roles = expandInheritedRole(model.InheritedRoleNames, databaseId.Name)
			}

			if err := client.MongoDBResourcesCreateUpdateMongoUserDefinitionThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CosmosDbMongoUserDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.MongoRBACClient

			id, err := mongorbacs.ParseMongodbUserDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.MongoDBResourcesGetMongoUserDefinition(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := CosmosDbMongoUserDefinitionResourceModel{}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					databaseId := parse.NewMongodbDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, *properties.DatabaseName)

					state.CosmosMongoDatabaseId = databaseId.ID()
					state.Username = pointer.From(properties.UserName)
					state.Password = metadata.ResourceData.Get("password").(string)
					state.InheritedRoleNames = flattenInheritedRole(properties.Roles)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CosmosDbMongoUserDefinitionResource) Delete() sdk.ResourceFunc {
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

			if err := client.MongoDBResourcesDeleteMongoUserDefinitionThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

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

func flattenInheritedRole(input *[]mongorbacs.Role) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, pointer.From(v.Role))
	}

	return result
}

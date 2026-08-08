// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CosmosDbMongoRoleDefinitionResourceModel struct {
	CosmosMongoDatabaseId string      `tfschema:"cosmos_mongo_database_id"`
	RoleName              string      `tfschema:"role_name"`
	InheritedRoleNames    []string    `tfschema:"inherited_role_names"`
	Privileges            []Privilege `tfschema:"privilege"`
}

type Privilege struct {
	Actions  []string   `tfschema:"actions"`
	Resource []Resource `tfschema:"resource"`
}

type Resource struct {
	CollectionName string `tfschema:"collection_name"`
	DbName         string `tfschema:"db_name"`
}

type CosmosDbMongoRoleDefinitionResource struct{}

var _ sdk.ResourceWithUpdate = CosmosDbMongoRoleDefinitionResource{}

func (r CosmosDbMongoRoleDefinitionResource) ResourceType() string {
	return "azurerm_cosmosdb_mongo_role_definition"
}

func (r CosmosDbMongoRoleDefinitionResource) ModelObject() interface{} {
	return &CosmosDbMongoRoleDefinitionResourceModel{}
}

func (r CosmosDbMongoRoleDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return openapis.ValidateMongodbRoleDefinitionID
}

func (r CosmosDbMongoRoleDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cosmos_mongo_database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: openapis.ValidateMongodbDatabaseID,
		},

		"role_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
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

		"privilege": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"resource": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"collection_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.CosmosEntityName,
								},

								"db_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.CosmosEntityName,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r CosmosDbMongoRoleDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CosmosDbMongoRoleDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbMongoRoleDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.OpenapisClient
			databaseId, err := openapis.ParseMongodbDatabaseID(model.CosmosMongoDatabaseId)
			if err != nil {
				return err
			}

			mongoRoleDefinitionId := fmt.Sprintf("%s.%s", databaseId.MongodbDatabaseName, model.RoleName)
			id := openapis.NewMongodbRoleDefinitionID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.DatabaseAccountName, mongoRoleDefinitionId)

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

			roleType := openapis.MongoRoleDefinitionTypeCustomRole
			parameters := openapis.MongoRoleDefinitionCreateUpdateParameters{
				Properties: &openapis.MongoRoleDefinitionResource{
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

func (r CosmosDbMongoRoleDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.OpenapisClient

			id, err := openapis.ParseMongodbRoleDefinitionID(metadata.ResourceData.Id())
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

			databaseId, err := openapis.ParseMongodbDatabaseID(model.CosmosMongoDatabaseId)
			if err != nil {
				return err
			}

			roleType := openapis.MongoRoleDefinitionTypeCustomRole
			parameters := openapis.MongoRoleDefinitionCreateUpdateParameters{
				Properties: &openapis.MongoRoleDefinitionResource{
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

func (r CosmosDbMongoRoleDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.OpenapisClient

			id, err := openapis.ParseMongodbRoleDefinitionID(metadata.ResourceData.Id())
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
					databaseId := openapis.NewMongodbDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, pointer.From(properties.DatabaseName))

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

func (r CosmosDbMongoRoleDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.OpenapisClient

			id, err := openapis.ParseMongodbRoleDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
			defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

			if err := client.MongoDBResourcesDeleteMongoRoleDefinitionThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandPrivilege(input []Privilege) *[]openapis.Privilege {
	if len(input) == 0 {
		return nil
	}

	result := make([]openapis.Privilege, 0, len(input))
	for _, v := range input {
		output := openapis.Privilege{
			Actions:  pointer.To(v.Actions),
			Resource: expandResource(v.Resource),
		}

		result = append(result, output)
	}

	return &result
}

func flattenPrivilege(input *[]openapis.Privilege) []Privilege {
	if input == nil {
		return []Privilege{}
	}

	result := make([]Privilege, 0, len(*input))
	for _, input := range *input {
		privilege := Privilege{
			Actions:  pointer.From(input.Actions),
			Resource: flattenResource(input.Resource),
		}

		result = append(result, privilege)
	}

	return result
}

func expandResource(input []Resource) *openapis.PrivilegeResource {
	if len(input) == 0 {
		return nil
	}

	privilegeResource := &input[0]
	result := openapis.PrivilegeResource{}

	if privilegeResource.CollectionName != "" {
		result.Collection = pointer.To(privilegeResource.CollectionName)
	}

	if privilegeResource.DbName != "" {
		result.Db = pointer.To(privilegeResource.DbName)
	}

	return &result
}

func flattenResource(input *openapis.PrivilegeResource) []Resource {
	var result []Resource
	if input == nil {
		return result
	}

	resource := Resource{
		CollectionName: pointer.From(input.Collection),
		DbName:         pointer.From(input.Db),
	}

	return append(result, resource)
}

func expandInheritedRoles(input []string, dbName string) *[]openapis.Role {
	if len(input) == 0 {
		return nil
	}

	result := make([]openapis.Role, 0, len(input))
	for _, v := range input {
		inheritedRole := openapis.Role{
			Db:   pointer.To(dbName),
			Role: pointer.To(v),
		}

		result = append(result, inheritedRole)
	}

	return &result
}

func flattenInheritedRoles(input *[]openapis.Role) []string {
	if input == nil {
		return []string{}
	}

	result := make([]string, 0, len(*input))
	for _, v := range *input {
		result = append(result, pointer.From(v.Role))
	}

	return result
}

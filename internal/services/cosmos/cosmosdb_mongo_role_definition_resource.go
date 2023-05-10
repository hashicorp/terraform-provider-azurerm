package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-11-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDbMongoRoleDefinitionResourceModel struct {
	AccountId          string      `tfschema:"account_id"`
	DbName             string      `tfschema:"db_name"`
	RoleName           string      `tfschema:"role_name"`
	InheritedRoleNames []string    `tfschema:"inherited_role_names"`
	Privileges         []Privilege `tfschema:"privilege"`
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
	return mongorbacs.ValidateMongodbRoleDefinitionID
}

func (r CosmosDbMongoRoleDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cosmosdb.ValidateDatabaseAccountID,
		},

		"db_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
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

			client := metadata.Client.Cosmos.MongoRBACClient
			databaseAccountId, err := cosmosdb.ParseDatabaseAccountID(model.AccountId)
			if err != nil {
				return err
			}

			mongoRoleDefinitionId := fmt.Sprintf("%s.%s", model.DbName, model.RoleName)
			id := mongorbacs.NewMongodbRoleDefinitionID(databaseAccountId.SubscriptionId, databaseAccountId.ResourceGroupName, databaseAccountId.DatabaseAccountName, mongoRoleDefinitionId)

			locks.ByName(id.DatabaseAccountName, CosmosDbAccountResourceName)
			defer locks.UnlockByName(id.DatabaseAccountName, CosmosDbAccountResourceName)

			existing, err := client.MongoDBResourcesGetMongoRoleDefinition(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			roleType := mongorbacs.MongoRoleDefinitionTypeOne
			properties := mongorbacs.MongoRoleDefinitionCreateUpdateParameters{
				Properties: &mongorbacs.MongoRoleDefinitionResource{
					DatabaseName: &model.DbName,
					RoleName:     &model.RoleName,
					Privileges:   expandPrivilege(model.Privileges),
					Roles:        expandInheritedRole(model.InheritedRoleNames, model.DbName),
					Type:         &roleType,
				},
			}

			if err := client.MongoDBResourcesCreateUpdateMongoRoleDefinitionThenPoll(ctx, id, properties); err != nil {
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

			roleType := mongorbacs.MongoRoleDefinitionTypeOne
			parameters := mongorbacs.MongoRoleDefinitionCreateUpdateParameters{
				Properties: &mongorbacs.MongoRoleDefinitionResource{
					DatabaseName: &model.DbName,
					RoleName:     &model.RoleName,
					Privileges:   expandPrivilege(model.Privileges),
					Roles:        expandInheritedRole(model.InheritedRoleNames, model.DbName),
					Type:         &roleType,
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

			state := CosmosDbMongoRoleDefinitionResourceModel{
				AccountId: cosmosdb.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName).ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.DbName = *properties.DatabaseName
					state.RoleName = *properties.RoleName
					state.Privileges = flattenPrivilege(properties.Privileges)
					state.InheritedRoleNames = flattenInheritedRole(properties.Roles)
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
			client := metadata.Client.Cosmos.MongoRBACClient

			id, err := mongorbacs.ParseMongodbRoleDefinitionID(metadata.ResourceData.Id())
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

func expandPrivilege(input []Privilege) *[]mongorbacs.Privilege {
	if len(input) == 0 {
		return nil
	}

	var result []mongorbacs.Privilege

	for _, v := range input {
		output := mongorbacs.Privilege{
			Actions:  &v.Actions,
			Resource: expandResource(v.Resource),
		}

		result = append(result, output)
	}

	return &result
}

func flattenPrivilege(input *[]mongorbacs.Privilege) []Privilege {
	var result []Privilege
	if input == nil {
		return result
	}

	for _, input := range *input {
		privilege := Privilege{
			Actions:  *input.Actions,
			Resource: flattenResource(input.Resource),
		}

		result = append(result, privilege)
	}

	return result
}

func expandResource(input []Resource) *mongorbacs.PrivilegeResource {
	if len(input) == 0 {
		return nil
	}

	privilegeResource := &input[0]
	result := mongorbacs.PrivilegeResource{}

	if privilegeResource.CollectionName != "" {
		result.Collection = &privilegeResource.CollectionName
	}

	if privilegeResource.DbName != "" {
		result.Db = &privilegeResource.DbName
	}

	return &result
}

func flattenResource(input *mongorbacs.PrivilegeResource) []Resource {
	var result []Resource
	if input == nil {
		return result
	}

	resource := Resource{}

	if input.Collection != nil {
		resource.CollectionName = *input.Collection
	}

	if input.Db != nil {
		resource.DbName = *input.Db
	}

	return append(result, resource)
}

func expandInheritedRole(input []string, dbName string) *[]mongorbacs.Role {
	if len(input) == 0 || dbName == "" {
		return nil
	}

	var result []mongorbacs.Role

	for _, v := range input {
		inheritedRole := mongorbacs.Role{
			Role: utils.String(v),
			Db:   utils.String(dbName),
		}

		result = append(result, inheritedRole)
	}

	return &result
}

func flattenInheritedRole(input *[]mongorbacs.Role) []string {
	var result []string
	if input == nil {
		return result
	}

	for _, input := range *input {
		var roleName string
		if role := input.Role; role != nil {
			roleName = *role
		}

		result = append(result, roleName)
	}

	return result
}

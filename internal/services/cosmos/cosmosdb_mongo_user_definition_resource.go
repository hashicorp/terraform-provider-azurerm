package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-11-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDbMongoUserDefinitionModel struct {
	AccountId  string `tfschema:"account_id"`
	CustomData string `tfschema:"custom_data"`
	DBName     string `tfschema:"db_name"`
	Password   string `tfschema:"password"`
	Roles      []Role `tfschema:"roles"`
	Username   string `tfschema:"username"`
}

type Role struct {
	Db   string `tfschema:"db"`
	Role string `tfschema:"role"`
}

type CosmosDbMongoUserDefinitionResource struct{}

var _ sdk.ResourceWithUpdate = CosmosDbMongoUserDefinitionResource{}

func (r CosmosDbMongoUserDefinitionResource) ResourceType() string {
	return "azurerm_cosmosdb_mongo_user_definition"
}

func (r CosmosDbMongoUserDefinitionResource) ModelObject() interface{} {
	return &CosmosDbMongoUserDefinitionModel{}
}

func (r CosmosDbMongoUserDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return mongorbacs.ValidateMongodbUserDefinitionID
}

func (r CosmosDbMongoUserDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: mongorbacs.ValidateDatabaseAccountID,
		},

		"db_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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

		"custom_data": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"roles": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"db": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"role": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
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
			var model CosmosDbMongoUserDefinitionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.MongoRBACClient
			databaseAccountId, err := mongorbacs.ParseDatabaseAccountID(model.AccountId)
			if err != nil {
				return err
			}

			mongoUserDefinitionId := fmt.Sprintf("%s.%s", model.DBName, model.Username)
			id := mongorbacs.NewMongodbUserDefinitionID(databaseAccountId.SubscriptionId, databaseAccountId.ResourceGroupName, databaseAccountId.DatabaseAccountName, mongoUserDefinitionId)

			existing, err := client.MongoDBResourcesGetMongoUserDefinition(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &mongorbacs.MongoUserDefinitionCreateUpdateParameters{
				Properties: &mongorbacs.MongoUserDefinitionResource{
					DatabaseName: &model.DBName,
					Mechanisms:   utils.String("SCRAM-SHA-256"),
					Password:     &model.Password,
					UserName:     &model.Username,
				},
			}

			if model.CustomData != "" {
				properties.Properties.CustomData = &model.CustomData
			}

			rolesValue, err := expandRole(model.Roles)
			if err != nil {
				return err
			}
			properties.Properties.Roles = rolesValue

			if err := client.MongoDBResourcesCreateUpdateMongoUserDefinitionThenPoll(ctx, id, *properties); err != nil {
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

			var model CosmosDbMongoUserDefinitionModel
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

			parameters := mongorbacs.MongoUserDefinitionCreateUpdateParameters{
				Properties: &mongorbacs.MongoUserDefinitionResource{
					DatabaseName: &model.DBName,
					Mechanisms:   utils.String("SCRAM-SHA-256"),
					Password:     &model.Password,
					UserName:     &model.Username,
				},
			}

			if model.CustomData != "" {
				properties.Properties.CustomData = &model.CustomData
			}

			rolesValue, err := expandRole(model.Roles)
			if err != nil {
				return err
			}
			properties.Properties.Roles = rolesValue

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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := CosmosDbMongoUserDefinitionModel{
				AccountId: mongorbacs.NewDatabaseAccountID(id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName).ID(),
			}

			if properties := model.Properties; properties != nil {
				state.DBName = *properties.DatabaseName
				state.Username = *properties.UserName
				state.Password = metadata.ResourceData.Get("password").(string)
				state.CustomData = metadata.ResourceData.Get("custom_data").(string)

				rolesValue, err := flattenRole(properties.Roles)
				if err != nil {
					return err
				}
				state.Roles = rolesValue
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

			if err := client.MongoDBResourcesDeleteMongoUserDefinitionThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandRole(input []Role) (*[]mongorbacs.Role, error) {
	var result []mongorbacs.Role

	for _, v := range input {
		input := v
		output := mongorbacs.Role{}

		if input.Db != "" {
			output.Db = &input.Db
		}

		if input.Role != "" {
			output.Role = &input.Role
		}

		result = append(result, output)
	}

	return &result, nil
}

func flattenRole(input *[]mongorbacs.Role) ([]Role, error) {
	var result []Role
	if input == nil {
		return result, nil
	}

	for _, input := range *input {
		output := Role{}

		if input.Db != nil {
			output.Db = *input.Db
		}

		if input.Role != nil {
			output.Role = *input.Role
		}

		result = append(result, output)
	}

	return result, nil
}

package postgresqlhsc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PostgreSQLHyperScaleRoleResourceModel struct {
	Name          string `tfschema:"name"`
	ServerGroupId string `tfschema:"server_group_id"`
	Password      string `tfschema:"password"`
}

type PostgreSQLHyperScaleRoleResource struct{}

var _ sdk.Resource = PostgreSQLHyperScaleRoleResource{}

func (r PostgreSQLHyperScaleRoleResource) ResourceType() string {
	return "azurerm_postgresql_hyperscale_role"
}

func (r PostgreSQLHyperScaleRoleResource) ModelObject() interface{} {
	return &PostgreSQLHyperScaleRoleResourceModel{}
}

func (r PostgreSQLHyperScaleRoleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return roles.ValidateRoleID
}

func (r PostgreSQLHyperScaleRoleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"server_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: roles.ValidateServerGroupsv2ID,
		},

		"password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r PostgreSQLHyperScaleRoleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PostgreSQLHyperScaleRoleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PostgreSQLHyperScaleRoleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PostgreSQLHSC.RolesClient
			serverGroupId, err := roles.ParseServerGroupsv2ID(model.ServerGroupId)
			if err != nil {
				return err
			}

			id := roles.NewRoleID(serverGroupId.SubscriptionId, serverGroupId.ResourceGroupName, serverGroupId.ServerGroupsv2Name, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := &roles.Role{
				Properties: roles.RoleProperties{
					Password: model.Password,
				},
			}

			if err := client.CreateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PostgreSQLHyperScaleRoleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PostgreSQLHSC.RolesClient

			id, err := roles.ParseRoleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
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

			state := PostgreSQLHyperScaleRoleResourceModel{
				Name:          id.RoleName,
				ServerGroupId: roles.NewServerGroupsv2ID(id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name).ID(),
			}

			properties := &model.Properties
			state.Password = properties.Password

			return metadata.Encode(&state)
		},
	}
}

func (r PostgreSQLHyperScaleRoleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PostgreSQLHSC.RolesClient

			id, err := roles.ParseRoleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

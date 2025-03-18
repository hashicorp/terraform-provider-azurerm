package apimanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceDataSourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	ApiManagementId   string `tfschema:"api_management_id"`
	DisplayName       string `tfschema:"display_name"`
}

type ApiManagementWorkspaceDataSource struct{}

var _ sdk.DataSource = ApiManagementWorkspaceDataSource{}

func (d ApiManagementWorkspaceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"api_management_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d ApiManagementWorkspaceDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d ApiManagementWorkspaceDataSource) ModelObject() interface{} {
	return &ApiManagementWorkspaceDataSourceModel{}
}

func (d ApiManagementWorkspaceDataSource) ResourceType() string {
	return "azurerm_api_management_workspace"
}

func (d ApiManagementWorkspaceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiManagementWorkspaceDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := apimanagementservice.ParseServiceID(model.ApiManagementId)
			if err != nil {
				return err
			}

			workspaceId := workspace.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.ServiceName, model.Name)

			resp, err := client.Get(ctx, workspaceId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", workspaceId)
				}
				return fmt.Errorf("retrieving %s: %+v", workspaceId, err)
			}

			state := ApiManagementWorkspaceDataSourceModel{
				Name:              workspaceId.WorkspaceId,
				ApiManagementId:   model.ApiManagementId,
				ResourceGroupName: workspaceId.ResourceGroup,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
				}
			}

			metadata.SetID(workspaceId)
			return metadata.Encode(&state)
		},
	}
}

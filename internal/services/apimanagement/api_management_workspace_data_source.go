package apimanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceDataSourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	ServiceName       string `tfschema:"service_name"`
	WorkspaceName     string `tfschema:"workspace_name"`
}

type ApiManagementWorkspaceDataSource struct{}

var _ sdk.DataSource = ApiManagementWorkspaceDataSource{}

func (d ApiManagementWorkspaceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"service_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d ApiManagementWorkspaceDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_name": {
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

			id := workspace.NewWorkspaceID(subscriptionId, model.ResourceGroupName, model.ServiceName, model.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ApiManagementWorkspaceDataSourceModel{
				Name:              id.WorkspaceId,
				ServiceName:       id.ServiceName,
				ResourceGroupName: id.ResourceGroup,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.WorkspaceName = props.DisplayName
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

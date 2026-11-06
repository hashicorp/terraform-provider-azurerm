package api_center_workspace

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *workspaces.WorkspaceId, model *workspaces.Workspace) error {
	state := ApiCenterWorkspaceModel{
		Name:               id.WorkspaceName,
		ApiCenterServiceID: services.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName).ID(),
	}

	if props := model.Properties; props != nil {
		state.Description = pointer.From(props.Description)
		state.Title = props.Title
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

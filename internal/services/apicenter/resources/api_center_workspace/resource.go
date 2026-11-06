package api_center_workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name api_center_workspace -service-package-name apicenter -properties "name" -compare-values "subscription_id:api_center_service_id,resource_group_name:api_center_service_id,service_name:api_center_service_id"

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.Resource             = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &workspaces.WorkspaceId{}
}

type Resource struct{}

func (r Resource) ModelObject() interface{} {
	return &ApiCenterWorkspaceModel{}
}

func (r Resource) ResourceType() string {
	return "azurerm_api_center_workspace"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspaces.ValidateWorkspaceID
}

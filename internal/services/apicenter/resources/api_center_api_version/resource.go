package api_center_api_version

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name azurerm_api_center_api_version -service-package-name apicenter -properties "name" -compare-values "subscription_id:api_id,resource_group_name:api_id,service_name:api_id,workspace_name:api_id,api_name:api_id"

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.Resource             = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &apiversions.VersionId{}
}

type Resource struct{}

func (r Resource) ModelObject() interface{} {
	return &ApiCenterApiVersionModel{}
}

func (r Resource) ResourceType() string {
	return "azurerm_api_center_api_version"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apiversions.ValidateVersionID
}

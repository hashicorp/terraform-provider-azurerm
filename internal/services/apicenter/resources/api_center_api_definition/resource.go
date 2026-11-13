// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name azurerm_api_center_api_definition -service-package-name apicenter -properties "name" -compare-values "subscription_id:api_version_id,resource_group_name:api_version_id,service_name:api_version_id,workspace_name:api_version_id,api_name:api_version_id,version_name:api_version_id"

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.Resource             = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &apidefinitions.DefinitionId{}
}

type Resource struct{}

func (r Resource) ModelObject() interface{} {
	return &ApiCenterApiDefinitionModel{}
}

func (r Resource) ResourceType() string {
	return "azurerm_api_center_api_definition"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apidefinitions.ValidateDefinitionID
}

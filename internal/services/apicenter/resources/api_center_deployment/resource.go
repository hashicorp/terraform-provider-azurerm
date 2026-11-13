// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_deployment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name azurerm_api_center_deployment -service-package-name apicenter -properties "name" -compare-values "subscription_id:api_definition_id,resource_group_name:api_definition_id,service_name:api_definition_id,workspace_name:api_definition_id,api_name:api_definition_id"

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.Resource             = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &deployments.DeploymentId{}
}

type Resource struct{}

func (r Resource) ModelObject() interface{} {
	return &ApiCenterDeploymentModel{}
}

func (r Resource) ResourceType() string {
	return "azurerm_api_center_deployment"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deployments.ValidateDeploymentID
}

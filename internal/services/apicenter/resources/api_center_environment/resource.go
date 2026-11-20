// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_environment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name azurerm_api_center_environment -service-package-name apicenter -properties "name" -compare-values "subscription_id:api_center_service_id,resource_group_name:api_center_service_id,service_name:api_center_service_id" -known-values "workspace_name:default"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &environments.EnvironmentId{}
}

func (r Resource) IdentityType() pluginsdk.ResourceTypeForIdentity {
	return pluginsdk.ResourceTypeForIdentityDefault
}

func (r Resource) ResourceType() string {
	return "azurerm_api_center_environment"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return environments.ValidateEnvironmentID
}

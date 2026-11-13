// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name azurerm_api_center_api -service-package-name apicenter -properties "name" -compare-values "subscription_id:service_id,resource_group_name:service_id,service_name:service_id" -known-values "workspace_name:default"

var _ sdk.ResourceWithIdentity = Resource{}

type Resource struct{}

func (r Resource) Identity() resourceids.ResourceId {
	return &apis.ApiId{}
}

func (r Resource) IdentityType() pluginsdk.ResourceTypeForIdentity {
	return pluginsdk.ResourceTypeForIdentityDefault
}

func (r Resource) ModelObject() any {
	return &ApiCenterApiModel{}
}

func (r Resource) ResourceType() string {
	return "azurerm_api_center_api"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apis.ValidateApiID
}

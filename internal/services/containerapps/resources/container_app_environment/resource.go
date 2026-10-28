// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithCustomizeDiff = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_container_app_environment"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedenvironments.ValidateManagedEnvironmentID
}

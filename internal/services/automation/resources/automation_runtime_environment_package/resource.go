// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment_package

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_automation_runtime_environment_package"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return packageresource.ValidatePackageID
}

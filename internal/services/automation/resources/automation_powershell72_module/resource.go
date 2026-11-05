// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_powershell72_module

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_automation_powershell72_module"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return module.ValidatePowerShell72ModuleID
}

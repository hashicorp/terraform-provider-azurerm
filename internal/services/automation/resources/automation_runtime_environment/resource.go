// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (m Resource) ResourceType() string {
	return "azurerm_automation_runtime_environment"
}

func (m Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return runtimeenvironment.ValidateRuntimeEnvironmentID
}

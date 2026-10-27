// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_deployment_script_azure_power_shell

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_resource_deployment_script_azure_power_shell"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deploymentscripts.ValidateDeploymentScriptID
}

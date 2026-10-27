// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_deployment_script_azure_cli

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return helpers.GetDeploymentScriptArguments(helpers.AzureCliKind)
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return helpers.GetDeploymentScriptAttributes()
}

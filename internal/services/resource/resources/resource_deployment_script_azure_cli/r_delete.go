// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_deployment_script_azure_cli

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/helpers"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return helpers.DeleteDeploymentScript()
}

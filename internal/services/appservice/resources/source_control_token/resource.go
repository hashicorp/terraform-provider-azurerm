// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package source_control_token

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate = Resource{}
	_ sdk.Resource           = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_source_control_token"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{"/providers/Microsoft.Web/sourceControls/GitHub"}, false)
}

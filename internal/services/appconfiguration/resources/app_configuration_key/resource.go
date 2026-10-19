// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_key

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithCustomizeDiff = Resource{}

var _ sdk.ResourceWithStateMigration = Resource{}

func (k Resource) ResourceType() string {
	return "azurerm_app_configuration_key"
}

func (k Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NestedItemId
}

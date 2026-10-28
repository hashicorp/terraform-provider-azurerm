// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_environment_type

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return environmenttypes.ValidateDevCenterEnvironmentTypeID
}

func (r Resource) ResourceType() string {
	return "azurerm_dev_center_environment_type"
}

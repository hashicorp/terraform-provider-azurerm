// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package aadb2c_directory

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_aadb2c_directory"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tenants.ValidateB2CDirectoryID
}

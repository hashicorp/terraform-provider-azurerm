// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_catalog

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return catalogs.ValidateDevCenterCatalogID
}

func (r Resource) ResourceType() string {
	return "azurerm_dev_center_catalog"
}

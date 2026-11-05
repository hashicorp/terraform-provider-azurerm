// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_catalog

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbitalplanetarycomputer/2026-04-15/geocatalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_geo_catalog"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return geocatalogs.ValidateGeoCatalogID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &geocatalogs.GeoCatalogId{}
}

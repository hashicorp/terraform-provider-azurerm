// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_bookshelf

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_discovery_bookshelf"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return bookshelves.ValidateBookshelfID
}

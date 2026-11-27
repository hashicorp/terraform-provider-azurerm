// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_shared_private_link_service

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_search_shared_private_link_service"
}

func (r Resource) Description() string {
	return "Manages the Shared Private Link Service for an Azure Search Service."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sharedprivatelinkresources.ValidateSharedPrivateLinkResourceID
}

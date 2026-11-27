// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service_datasource_blob

import (
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_search_service_datasource_blob"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datasources.ValidateDatasourceID
}

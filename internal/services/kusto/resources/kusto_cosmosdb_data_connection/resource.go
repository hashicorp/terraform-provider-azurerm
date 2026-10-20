// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cosmosdb_data_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_kusto_cosmosdb_data_connection"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataconnections.ValidateDataConnectionID
}

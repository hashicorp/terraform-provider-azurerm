// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_dedicated_gateway

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_cosmosdb_sql_dedicated_gateway"
}

func (r Resource) Description() string {
	return "Manages a SQL Dedicated Gateway within a Cosmos DB Account."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sqldedicatedgateway.ValidateServiceID
}

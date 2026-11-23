// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_role_definition

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var CosmosDbAccountResourceName = "azurerm_cosmosdb_account"

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_cosmosdb_mongo_role_definition"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return mongorbacs.ValidateMongodbRoleDefinitionID
}

// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_node_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var CosmosDbPostgreSQLClusterResourceName = "azurerm_cosmosdb_postgresql_cluster"

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_cosmosdb_postgresql_node_configuration"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurations.ValidateNodeConfigurationID
}

// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_coordinator_configuration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name cosmosdb_postgresql_coordinator_configuration -properties "name" -compare-values "subscription_id:cluster_id,resource_group_name:cluster_id,server_groups_v2_name:cluster_id" -test-params "on"

var CosmosDbPostgreSQLClusterResourceName = "azurerm_cosmosdb_postgresql_cluster"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_cosmosdb_postgresql_coordinator_configuration"
}

func (r Resource) Description() string {
	return "Sets a Coordinator Configuration value on Azure Cosmos DB for PostgreSQL Cluster."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurations.ValidateCoordinatorConfigurationID
}

func (r Resource) Identity() resourceids.ResourceId {
	return new(configurations.CoordinatorConfigurationId)
}

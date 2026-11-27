// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name cosmosdb_postgresql_cluster -properties "name,resource_group_name" -test-expect-non-empty

var CosmosDbPostgreSQLClusterResourceName = "azurerm_cosmosdb_postgresql_cluster"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (r Resource) ResourceType() string {
	return CosmosDbPostgreSQLClusterResourceName
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return clusters.ValidateServerGroupsv2ID
}

func (r Resource) Identity() resourceids.ResourceId {
	return new(clusters.ServerGroupsv2Id)
}

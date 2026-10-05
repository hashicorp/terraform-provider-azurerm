// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Azure Managed Redis (AMR) consists of two ARM resource types: cluster and database. Cluster is where compute, load
// balancer, network and other infrastructure is setup. Database refers to the Redis instance / process itself. Database
// is a child of cluster with 1-1 mapping.
//
// Database was its own resource in the deprecated redis_enterprise resource, but intentionally included here to improve
// UX. There were cases where users not aware of database and expect Redis Enterprise cluster to work by itself.
//
// There might be a plan to support multiple databases in the future, in which case we will implement an
// azurerm_managed_redis_custom_db resource.

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_managed_redis"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return redisenterprise.ValidateRedisEnterpriseID
}

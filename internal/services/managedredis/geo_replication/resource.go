// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_replication

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Geo-replication / linked databases are managed as a separate resource because when dbs are linked, ARM
// will mutate the state of all dbs out of bound, causing unexpected plan diff.
//
// The default database name is always "default", and because cluster and database are managed as a single TF resource,
// we use cluster id to configure linking. Internally the database id can be derived from cluster id
// by appending "/databases/default" suffix.

type Resource struct{}

var (
	_ sdk.ResourceWithCustomizeDiff = Resource{}
	_ sdk.ResourceWithUpdate        = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_managed_redis_geo_replication"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return redisenterprise.ValidateRedisEnterpriseID
}

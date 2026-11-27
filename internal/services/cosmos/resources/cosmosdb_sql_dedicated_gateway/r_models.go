// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_dedicated_gateway

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
)

type CosmosDbSqlDedicatedGatewayModel struct {
	CosmosDbAccountId string                          `tfschema:"cosmosdb_account_id"`
	InstanceCount     int64                           `tfschema:"instance_count"`
	InstanceSize      sqldedicatedgateway.ServiceSize `tfschema:"instance_size"`
}

func (r Resource) ModelObject() any {
	return &CosmosDbSqlDedicatedGatewayModel{}
}

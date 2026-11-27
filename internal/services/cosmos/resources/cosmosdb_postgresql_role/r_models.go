// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_role

type CosmosDbPostgreSQLRoleResourceModel struct {
	Name      string `tfschema:"name"`
	ClusterId string `tfschema:"cluster_id"`
	Password  string `tfschema:"password"`
}

func (r Resource) ModelObject() any {
	return &CosmosDbPostgreSQLRoleResourceModel{}
}

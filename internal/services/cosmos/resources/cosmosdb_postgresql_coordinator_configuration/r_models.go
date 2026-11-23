// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_coordinator_configuration

type CosmosDbPostgreSQLCoordinatorConfigurationModel struct {
	Name      string `tfschema:"name"`
	ClusterId string `tfschema:"cluster_id"`
	Value     string `tfschema:"value"`
}

func (r Resource) ModelObject() any {
	return &Resource{}
}

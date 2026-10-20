// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_cosmosdb_data_connection

type CosmosDBDataConnectionModel struct {
	Name                string `tfschema:"name"`
	Location            string `tfschema:"location"`
	CosmosDbContainerId string `tfschema:"cosmosdb_container_id"`
	DatabaseId          string `tfschema:"kusto_database_id"`
	ManagedIdentityId   string `tfschema:"managed_identity_id"`
	MappingRuleName     string `tfschema:"mapping_rule_name"`
	RetrievalStartDate  string `tfschema:"retrieval_start_date"`
	TableName           string `tfschema:"table_name"`
}

func (r Resource) ModelObject() any {
	return &CosmosDBDataConnectionModel{}
}

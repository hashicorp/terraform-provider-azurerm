// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_time_series_database_connection

type TimeSeriesDatabaseConnectionModel struct {
	Name                         string `tfschema:"name"`
	DigitalTwinsId               string `tfschema:"digital_twins_id"`
	EventhubConsumerGroupName    string `tfschema:"eventhub_consumer_group_name"`
	EventhubName                 string `tfschema:"eventhub_name"`
	EventhubNamespaceEndpointUri string `tfschema:"eventhub_namespace_endpoint_uri"`
	EventhubNamespaceId          string `tfschema:"eventhub_namespace_id"`
	KustoClusterId               string `tfschema:"kusto_cluster_id"`
	KustoClusterUri              string `tfschema:"kusto_cluster_uri"`
	KustoDatabaseName            string `tfschema:"kusto_database_name"`
	KustoTableName               string `tfschema:"kusto_table_name"`
}

func (m Resource) ModelObject() any {
	return &TimeSeriesDatabaseConnectionModel{}
}

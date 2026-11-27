// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_cosmosdb

type OutputCosmosDBResourceModel struct {
	Name               string `tfschema:"name"`
	StreamAnalyticsJob string `tfschema:"stream_analytics_job_id"`
	AccountKey         string `tfschema:"cosmosdb_account_key"`
	Database           string `tfschema:"cosmosdb_sql_database_id"`
	ContainerName      string `tfschema:"container_name"`
	DocumentID         string `tfschema:"document_id"`
	PartitionKey       string `tfschema:"partition_key"`
	AuthenticationMode string `tfschema:"authentication_mode"`
}

func (r Resource) ModelObject() any {
	return &OutputCosmosDBResourceModel{}
}

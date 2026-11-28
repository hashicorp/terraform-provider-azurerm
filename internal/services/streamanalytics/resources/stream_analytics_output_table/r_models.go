// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_table

type OutputTableResourceModel struct {
	Name               string   `tfschema:"name"`
	StreamAnalyticsJob string   `tfschema:"stream_analytics_job_name"`
	ResourceGroup      string   `tfschema:"resource_group_name"`
	StorageAccount     string   `tfschema:"storage_account_name"`
	StorageAccountKey  string   `tfschema:"storage_account_key"`
	Table              string   `tfschema:"table"`
	PartitionKey       string   `tfschema:"partition_key"`
	RowKey             string   `tfschema:"row_key"`
	BatchSize          int64    `tfschema:"batch_size"`
	ColumnsToRemove    []string `tfschema:"columns_to_remove"`
}

func (r Resource) ModelObject() any {
	return &OutputTableResourceModel{}
}

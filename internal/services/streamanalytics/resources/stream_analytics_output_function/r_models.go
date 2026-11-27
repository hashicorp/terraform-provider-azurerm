// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_function

type OutputFunctionResourceModel struct {
	Name               string `tfschema:"name"`
	StreamAnalyticsJob string `tfschema:"stream_analytics_job_name"`
	ResourceGroup      string `tfschema:"resource_group_name"`
	FunctionApp        string `tfschema:"function_app"`
	FunctionName       string `tfschema:"function_name"`
	ApiKey             string `tfschema:"api_key"`
	BatchMaxInBytes    int64  `tfschema:"batch_max_in_bytes"`
	BatchMaxCount      int64  `tfschema:"batch_max_count"`
}

func (r Resource) ModelObject() any {
	return &OutputFunctionResourceModel{}
}

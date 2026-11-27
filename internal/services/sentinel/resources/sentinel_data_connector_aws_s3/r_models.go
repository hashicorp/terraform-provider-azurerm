// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_aws_s3

type DataConnectorAwsS3Model struct {
	Name                    string   `tfschema:"name"`
	LogAnalyticsWorkspaceId string   `tfschema:"log_analytics_workspace_id"`
	AwsRoleArm              string   `tfschema:"aws_role_arn"`
	DestinationTable        string   `tfschema:"destination_table"`
	SqsUrls                 []string `tfschema:"sqs_urls"`
}

func (r Resource) ModelObject() any {
	return &DataConnectorAwsS3Model{}
}

// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_powerbi

type OutputPowerBIResourceModel struct {
	Name                   string `tfschema:"name"`
	StreamAnalyticsJob     string `tfschema:"stream_analytics_job_id"`
	DataSet                string `tfschema:"dataset"`
	Table                  string `tfschema:"table"`
	GroupID                string `tfschema:"group_id"`
	GroupName              string `tfschema:"group_name"`
	TokenUserPrincipalName string `tfschema:"token_user_principal_name"`
	TokenUserDisplayName   string `tfschema:"token_user_display_name"`
}

func (r Resource) ModelObject() any {
	return &OutputPowerBIResourceModel{}
}

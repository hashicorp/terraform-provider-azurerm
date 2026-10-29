// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table

type LogAnalyticsWorkspaceTableResourceModel struct {
	Name                 string `tfschema:"name"`
	WorkspaceId          string `tfschema:"workspace_id"`
	Plan                 string `tfschema:"plan"`
	RetentionInDays      int64  `tfschema:"retention_in_days"`
	TotalRetentionInDays int64  `tfschema:"total_retention_in_days"`
}

func (r Resource) ModelObject() any {
	return &LogAnalyticsWorkspaceTableResourceModel{}
}

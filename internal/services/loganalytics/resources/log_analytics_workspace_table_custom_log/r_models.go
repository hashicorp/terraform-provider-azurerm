// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_custom_log

type WorkspaceTableCustomLogResourceModel struct {
	Name                 string                 `tfschema:"name"`
	WorkspaceId          string                 `tfschema:"workspace_id"`
	DisplayName          string                 `tfschema:"display_name"`
	Description          string                 `tfschema:"description"`
	Plan                 string                 `tfschema:"plan"`
	Columns              []workspaceTableColumn `tfschema:"column"`
	Solutions            []string               `tfschema:"solutions"`
	StandardColumns      []workspaceTableColumn `tfschema:"standard_column"`
	RetentionInDays      int64                  `tfschema:"retention_in_days"`
	TotalRetentionInDays int64                  `tfschema:"total_retention_in_days"`
}

func (r Resource) ModelObject() any {
	return &WorkspaceTableCustomLogResourceModel{}
}

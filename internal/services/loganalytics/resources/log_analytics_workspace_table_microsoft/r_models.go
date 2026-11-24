// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_microsoft

type WorkspaceTableMicrosoftResourceModel struct {
	Name                 string   `tfschema:"name"`
	WorkspaceId          string   `tfschema:"workspace_id"`
	DisplayName          string   `tfschema:"display_name"`
	Description          string   `tfschema:"description"`
	Columns              []Column `tfschema:"column"`
	Labels               []string `tfschema:"labels"`
	Solutions            []string `tfschema:"solutions"`
	StandardColumns      []Column `tfschema:"standard_column"`
	RetentionInDays      int64    `tfschema:"retention_in_days"`
	TotalRetentionInDays int64    `tfschema:"total_retention_in_days"`
}

func (r Resource) ModelObject() any {
	return &WorkspaceTableMicrosoftResourceModel{}
}

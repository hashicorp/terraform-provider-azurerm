// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_group

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspaceGroupModel{}
}

type ApiManagementWorkspaceGroupModel struct {
	Name                     string `tfschema:"name"`
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	DisplayName              string `tfschema:"display_name"`
	Description              string `tfschema:"description"`
	ExternalId               string `tfschema:"external_id"`
	Type                     string `tfschema:"type"`
}

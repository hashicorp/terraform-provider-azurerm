// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_api_version_set

type ApiManagementWorkspaceApiVersionSetModel struct {
	Name                     string `tfschema:"name"`
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	DisplayName              string `tfschema:"display_name"`
	VersioningScheme         string `tfschema:"versioning_scheme"`
	Description              string `tfschema:"description"`
	VersionHeaderName        string `tfschema:"version_header_name"`
	VersionQueryName         string `tfschema:"version_query_name"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspaceApiVersionSetModel{}
}

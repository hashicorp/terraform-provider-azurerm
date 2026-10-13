// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_policy

type ApiManagementWorkspacePolicyModel struct {
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	XmlContent               string `tfschema:"xml_content"`
	XmlLink                  string `tfschema:"xml_link"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspacePolicyModel{}
}

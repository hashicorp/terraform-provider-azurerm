// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_policy_fragment

type ApiManagementWorkspacePolicyFragmentModel struct {
	Name                     string `tfschema:"name"`
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	Description              string `tfschema:"description"`
	XmlFormat                string `tfschema:"xml_format"`
	XmlContent               string `tfschema:"xml_content"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspacePolicyFragmentModel{}
}

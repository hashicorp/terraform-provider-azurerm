// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_named_value

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspaceNamedValueModel{}
}

type ApiManagementWorkspaceNamedValueModel struct {
	ApiManagementWorkspaceId string     `tfschema:"api_management_workspace_id"`
	DisplayName              string     `tfschema:"display_name"`
	Name                     string     `tfschema:"name"`
	Secret                   bool       `tfschema:"secret"`
	Tags                     []string   `tfschema:"tags"`
	Value                    string     `tfschema:"value"`
	ValueFromKeyVault        []KeyVault `tfschema:"value_from_key_vault"`
}

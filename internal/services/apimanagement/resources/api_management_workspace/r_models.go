// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace

type ApiManagementWorkspaceModel struct {
	Name            string `tfschema:"name"`
	ApiManagementId string `tfschema:"api_management_id"`
	Description     string `tfschema:"description"`
	DisplayName     string `tfschema:"display_name"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspaceModel{}
}

// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_config_connection

type ApiManagementGatewayConfigConnectionModel struct {
	Name                     string   `tfschema:"name"`
	ApiManagementGatewayId   string   `tfschema:"api_management_gateway_id"`
	ApiManagementWorkspaceId string   `tfschema:"api_management_workspace_id"`
	Hostnames                []string `tfschema:"hostnames"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementGatewayConfigConnectionModel{}
}

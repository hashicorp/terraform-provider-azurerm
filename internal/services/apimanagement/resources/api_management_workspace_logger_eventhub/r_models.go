// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_logger_eventhub

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspaceLoggerEventhubModel{}
}

type ApiManagementWorkspaceLoggerEventhubModel struct {
	Name                     string          `tfschema:"name"`
	ApiManagementWorkspaceId string          `tfschema:"api_management_workspace_id"`
	Eventhub                 []EventhubModel `tfschema:"eventhub"`
	BufferingEnabled         bool            `tfschema:"buffering_enabled"`
	Description              string          `tfschema:"description"`
	ResourceId               string          `tfschema:"resource_id"`
}

type EventhubModel struct {
	Name                         string `tfschema:"name"`
	ConnectionString             string `tfschema:"connection_string"`
	EndpointUri                  string `tfschema:"endpoint_uri"`
	UserAssignedIdentityClientId string `tfschema:"user_assigned_identity_client_id"`
}

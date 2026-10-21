// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_scope_connection

type ManagerScopeConnectionModel struct {
	Name             string `tfschema:"name"`
	NetworkManagerId string `tfschema:"network_manager_id"`
	ConnectionState  string `tfschema:"connection_state"`
	Description      string `tfschema:"description"`
	ResourceId       string `tfschema:"target_scope_id"`
	TenantId         string `tfschema:"tenant_id"`
}

func (r Resource) ModelObject() any {
	return &ManagerScopeConnectionModel{}
}

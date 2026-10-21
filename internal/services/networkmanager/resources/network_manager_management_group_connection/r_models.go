// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_management_group_connection

type ManagerManagementGroupConnectionModel struct {
	Name              string `tfschema:"name"`
	ManagementGroupId string `tfschema:"management_group_id"`
	ConnectionState   string `tfschema:"connection_state"`
	Description       string `tfschema:"description"`
	NetworkManagerId  string `tfschema:"network_manager_id"`
}

func (r Resource) ModelObject() any {
	return &ManagerManagementGroupConnectionModel{}
}

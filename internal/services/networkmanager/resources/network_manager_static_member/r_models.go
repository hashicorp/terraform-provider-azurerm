// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_static_member

type ManagerStaticMemberModel struct {
	Name           string `tfschema:"name"`
	NetworkGroupId string `tfschema:"network_group_id"`
	TargetVNetId   string `tfschema:"target_virtual_network_id"`
	Region         string `tfschema:"region"`
}

func (r Resource) ModelObject() any {
	return &ManagerStaticMemberModel{}
}

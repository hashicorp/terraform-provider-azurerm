// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_network_group

type ManagerNetworkGroupModel struct {
	Name             string `tfschema:"name"`
	NetworkManagerId string `tfschema:"network_manager_id"`
	Description      string `tfschema:"description"`
	MemberType       string `tfschema:"member_type"`
}

func (r Resource) ModelObject() any {
	return &ManagerNetworkGroupModel{}
}

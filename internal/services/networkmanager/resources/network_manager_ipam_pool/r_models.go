// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool

func (Resource) ModelObject() any {
	return &ManagerIpamPoolResourceModel{}
}

type ManagerIpamPoolResourceModel struct {
	AddressPrefixes  []string          `tfschema:"address_prefixes"`
	Description      string            `tfschema:"description"`
	DisplayName      string            `tfschema:"display_name"`
	Location         string            `tfschema:"location"`
	Name             string            `tfschema:"name"`
	NetworkManagerId string            `tfschema:"network_manager_id"`
	ParentPoolName   string            `tfschema:"parent_pool_name"`
	Tags             map[string]string `tfschema:"tags"`
}

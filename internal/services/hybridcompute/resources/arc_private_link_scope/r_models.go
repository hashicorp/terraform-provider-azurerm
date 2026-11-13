// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_private_link_scope

type PrivateLinkScopeModel struct {
	Name                       string            `tfschema:"name"`
	Location                   string            `tfschema:"location"`
	ResourceGroupName          string            `tfschema:"resource_group_name"`
	Tags                       map[string]string `tfschema:"tags"`
	PublicNetworkAccessEnabled bool              `tfschema:"public_network_access_enabled"`
}

func (a Resource) ModelObject() any {
	return &PrivateLinkScopeModel{}
}

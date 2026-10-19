// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_network_connection

func (r Resource) ModelObject() any {
	return &DevCenterNetworkConnectionResourceModel{}
}

type DevCenterNetworkConnectionResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	DomainJoinType    string            `tfschema:"domain_join_type"`
	SubnetId          string            `tfschema:"subnet_id"`
	DomainName        string            `tfschema:"domain_name"`
	DomainPassword    string            `tfschema:"domain_password"`
	DomainUsername    string            `tfschema:"domain_username"`
	OrganizationUnit  string            `tfschema:"organization_unit"`
	Tags              map[string]string `tfschema:"tags"`
}

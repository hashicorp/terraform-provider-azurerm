// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool_address

type BackendAddressPoolAddressModel struct {
	Name                    string                      `tfschema:"name"`
	BackendAddressPoolId    string                      `tfschema:"backend_address_pool_id"`
	VirtualNetworkId        string                      `tfschema:"virtual_network_id"`
	IPAddress               string                      `tfschema:"ip_address"`
	FrontendIPConfiguration string                      `tfschema:"backend_address_ip_configuration_id"`
	PortMapping             []inboundNATRulePortMapping `tfschema:"inbound_nat_rule_port_mapping"`
}

func (r Resource) ModelObject() any {
	return &BackendAddressPoolAddressModel{}
}

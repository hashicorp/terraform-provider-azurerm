// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_virtual_network_appliance

type NetworkVirtualApplianceResourceModel struct {
	Name         string `tfschema:"name"`
	VirtualHubID string `tfschema:"virtual_hub_id"`
}

func (r Resource) ModelObject() any {
	return &NetworkVirtualApplianceResourceModel{}
}

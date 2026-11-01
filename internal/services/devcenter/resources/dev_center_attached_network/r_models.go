// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_attached_network

func (r Resource) ModelObject() any {
	return &DevCenterAttachedNetworkResourceModel{}
}

type DevCenterAttachedNetworkResourceModel struct {
	Name                string `tfschema:"name"`
	DevCenterId         string `tfschema:"dev_center_id"`
	NetworkConnectionId string `tfschema:"network_connection_id"`
}

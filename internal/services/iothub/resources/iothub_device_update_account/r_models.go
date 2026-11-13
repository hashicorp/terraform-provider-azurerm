// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_device_update_account

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
)

type IotHubDeviceUpdateAccountModel struct {
	Name                       string            `tfschema:"name"`
	ResourceGroupName          string            `tfschema:"resource_group_name"`
	Location                   string            `tfschema:"location"`
	HostName                   string            `tfschema:"host_name"`
	PublicNetworkAccessEnabled bool              `tfschema:"public_network_access_enabled"`
	Sku                        deviceupdates.SKU `tfschema:"sku"`
	Tags                       map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &IotHubDeviceUpdateAccountModel{}
}

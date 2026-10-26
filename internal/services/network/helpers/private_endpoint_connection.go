// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
)

func FlattenNetworkInterface(networkInterfaceId string) any {
	id, err := commonids.ParseNetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"id":   id.ID(),
			"name": id.NetworkInterfaceName,
		},
	}
}

func GetPrivateIpAddress(ctx context.Context, client *networkinterfaces.NetworkInterfacesClient, networkInterfaceId string) string {
	privateIpAddress := ""
	id, err := commonids.ParseNetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return privateIpAddress
	}

	resp, err := client.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return privateIpAddress
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if configs := props.IPConfigurations; configs != nil {
				for i, config := range *configs {
					if propFmt := config.Properties; propFmt != nil {
						if propFmt.PrivateIPAddress != nil && *propFmt.PrivateIPAddress != "" && i == 0 {
							privateIpAddress = *propFmt.PrivateIPAddress
						}
						break
					}
				}
			}
		}
	}

	return privateIpAddress
}

// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
)

func localNetworkGatewayFromId(localNetworkGatewayId string) (string, error) {
	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(localNetworkGatewayId)
	if err != nil {
		return "", err
	}

	return id.LocalNetworkGatewayName, nil
}

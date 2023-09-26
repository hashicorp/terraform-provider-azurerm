package apimanagementservice

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdditionalLocation struct {
	DisableGateway              *bool                             `json:"disableGateway,omitempty"`
	GatewayRegionalUrl          *string                           `json:"gatewayRegionalUrl,omitempty"`
	Location                    string                            `json:"location"`
	PlatformVersion             *PlatformVersion                  `json:"platformVersion,omitempty"`
	PrivateIPAddresses          *[]string                         `json:"privateIPAddresses,omitempty"`
	PublicIPAddressId           *string                           `json:"publicIpAddressId,omitempty"`
	PublicIPAddresses           *[]string                         `json:"publicIPAddresses,omitempty"`
	Sku                         ApiManagementServiceSkuProperties `json:"sku"`
	VirtualNetworkConfiguration *VirtualNetworkConfiguration      `json:"virtualNetworkConfiguration,omitempty"`
	Zones                       *zones.Schema                     `json:"zones,omitempty"`
}

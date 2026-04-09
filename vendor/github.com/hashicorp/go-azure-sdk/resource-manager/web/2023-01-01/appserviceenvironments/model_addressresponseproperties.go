package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddressResponseProperties struct {
	InternalIPAddress   *string             `json:"internalIpAddress,omitempty"`
	OutboundIPAddresses *[]string           `json:"outboundIpAddresses,omitempty"`
	ServiceIPAddress    *string             `json:"serviceIpAddress,omitempty"`
	VipMappings         *[]VirtualIPMapping `json:"vipMappings,omitempty"`
}

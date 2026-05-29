package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfile struct {
	AppNetworkResourceGroup            *string                    `json:"appNetworkResourceGroup,omitempty"`
	AppSubnetId                        *string                    `json:"appSubnetId,omitempty"`
	IngressConfig                      *IngressConfig             `json:"ingressConfig,omitempty"`
	OutboundIPs                        *NetworkProfileOutboundIPs `json:"outboundIPs,omitempty"`
	OutboundType                       *string                    `json:"outboundType,omitempty"`
	RequiredTraffics                   *[]RequiredTraffic         `json:"requiredTraffics,omitempty"`
	ServiceCidr                        *string                    `json:"serviceCidr,omitempty"`
	ServiceRuntimeNetworkResourceGroup *string                    `json:"serviceRuntimeNetworkResourceGroup,omitempty"`
	ServiceRuntimeSubnetId             *string                    `json:"serviceRuntimeSubnetId,omitempty"`
}

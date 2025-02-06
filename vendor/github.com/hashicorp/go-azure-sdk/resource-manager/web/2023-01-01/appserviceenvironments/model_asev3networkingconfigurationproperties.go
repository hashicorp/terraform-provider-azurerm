package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AseV3NetworkingConfigurationProperties struct {
	AllowNewPrivateEndpointConnections *bool     `json:"allowNewPrivateEndpointConnections,omitempty"`
	ExternalInboundIPAddresses         *[]string `json:"externalInboundIpAddresses,omitempty"`
	FtpEnabled                         *bool     `json:"ftpEnabled,omitempty"`
	InboundIPAddressOverride           *string   `json:"inboundIpAddressOverride,omitempty"`
	InternalInboundIPAddresses         *[]string `json:"internalInboundIpAddresses,omitempty"`
	LinuxOutboundIPAddresses           *[]string `json:"linuxOutboundIpAddresses,omitempty"`
	RemoteDebugEnabled                 *bool     `json:"remoteDebugEnabled,omitempty"`
	WindowsOutboundIPAddresses         *[]string `json:"windowsOutboundIpAddresses,omitempty"`
}

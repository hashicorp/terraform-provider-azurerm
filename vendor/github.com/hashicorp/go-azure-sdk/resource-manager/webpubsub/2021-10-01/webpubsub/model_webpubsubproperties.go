package webpubsub

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebPubSubProperties struct {
	DisableAadAuth             *bool                        `json:"disableAadAuth,omitempty"`
	DisableLocalAuth           *bool                        `json:"disableLocalAuth,omitempty"`
	ExternalIP                 *string                      `json:"externalIP,omitempty"`
	HostName                   *string                      `json:"hostName,omitempty"`
	HostNamePrefix             *string                      `json:"hostNamePrefix,omitempty"`
	LiveTraceConfiguration     *LiveTraceConfiguration      `json:"liveTraceConfiguration,omitempty"`
	NetworkACLs                *WebPubSubNetworkACLs        `json:"networkACLs,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *string                      `json:"publicNetworkAccess,omitempty"`
	PublicPort                 *int64                       `json:"publicPort,omitempty"`
	ResourceLogConfiguration   *ResourceLogConfiguration    `json:"resourceLogConfiguration,omitempty"`
	ServerPort                 *int64                       `json:"serverPort,omitempty"`
	SharedPrivateLinkResources *[]SharedPrivateLinkResource `json:"sharedPrivateLinkResources,omitempty"`
	Tls                        *WebPubSubTlsSettings        `json:"tls,omitempty"`
	Version                    *string                      `json:"version,omitempty"`
}

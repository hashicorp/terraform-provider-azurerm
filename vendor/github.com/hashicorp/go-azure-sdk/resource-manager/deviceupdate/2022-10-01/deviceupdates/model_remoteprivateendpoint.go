package deviceupdates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemotePrivateEndpoint struct {
	ConnectionDetails                   *[]ConnectionDetails            `json:"connectionDetails,omitempty"`
	Id                                  *string                         `json:"id,omitempty"`
	ImmutableResourceId                 *string                         `json:"immutableResourceId,omitempty"`
	ImmutableSubscriptionId             *string                         `json:"immutableSubscriptionId,omitempty"`
	Location                            *string                         `json:"location,omitempty"`
	ManualPrivateLinkServiceConnections *[]PrivateLinkServiceConnection `json:"manualPrivateLinkServiceConnections,omitempty"`
	PrivateLinkServiceConnections       *[]PrivateLinkServiceConnection `json:"privateLinkServiceConnections,omitempty"`
	PrivateLinkServiceProxies           *[]PrivateLinkServiceProxy      `json:"privateLinkServiceProxies,omitempty"`
	VnetTrafficTag                      *string                         `json:"vnetTrafficTag,omitempty"`
}

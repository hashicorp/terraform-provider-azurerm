package deviceupdates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceProxy struct {
	GroupConnectivityInformation            *[]GroupConnectivityInformation    `json:"groupConnectivityInformation,omitempty"`
	Id                                      *string                            `json:"id,omitempty"`
	RemotePrivateEndpointConnection         *RemotePrivateEndpointConnection   `json:"remotePrivateEndpointConnection,omitempty"`
	RemotePrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState `json:"remotePrivateLinkServiceConnectionState,omitempty"`
}

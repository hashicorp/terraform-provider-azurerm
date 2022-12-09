package packetcorecontrolplane

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCoreControlPlaneClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPacketCoreControlPlaneClientWithBaseURI(endpoint string) PacketCoreControlPlaneClient {
	return PacketCoreControlPlaneClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package packetcoredataplane

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCoreDataPlaneClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPacketCoreDataPlaneClientWithBaseURI(endpoint string) PacketCoreDataPlaneClient {
	return PacketCoreDataPlaneClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

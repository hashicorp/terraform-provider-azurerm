package virtualnetworks

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworksClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVirtualNetworksClientWithBaseURI(endpoint string) VirtualNetworksClient {
	return VirtualNetworksClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

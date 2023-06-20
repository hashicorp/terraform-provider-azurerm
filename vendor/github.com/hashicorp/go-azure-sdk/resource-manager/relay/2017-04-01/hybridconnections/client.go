package hybridconnections

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewHybridConnectionsClientWithBaseURI(endpoint string) HybridConnectionsClient {
	return HybridConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

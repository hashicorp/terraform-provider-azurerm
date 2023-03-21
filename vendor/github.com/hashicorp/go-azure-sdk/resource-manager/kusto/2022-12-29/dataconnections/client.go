package dataconnections

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDataConnectionsClientWithBaseURI(endpoint string) DataConnectionsClient {
	return DataConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

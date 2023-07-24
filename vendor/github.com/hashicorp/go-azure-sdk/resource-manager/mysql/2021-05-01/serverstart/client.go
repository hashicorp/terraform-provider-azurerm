package serverstart

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerStartClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServerStartClientWithBaseURI(endpoint string) ServerStartClient {
	return ServerStartClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

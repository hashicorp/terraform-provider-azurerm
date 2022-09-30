package serverrestart

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerRestartClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServerRestartClientWithBaseURI(endpoint string) ServerRestartClient {
	return ServerRestartClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

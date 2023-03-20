package databases

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabasesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDatabasesClientWithBaseURI(endpoint string) DatabasesClient {
	return DatabasesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

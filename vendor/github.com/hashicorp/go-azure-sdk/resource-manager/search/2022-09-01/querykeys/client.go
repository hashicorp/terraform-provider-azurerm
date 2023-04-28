package querykeys

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryKeysClient struct {
	Client  autorest.Client
	baseUri string
}

func NewQueryKeysClientWithBaseURI(endpoint string) QueryKeysClient {
	return QueryKeysClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

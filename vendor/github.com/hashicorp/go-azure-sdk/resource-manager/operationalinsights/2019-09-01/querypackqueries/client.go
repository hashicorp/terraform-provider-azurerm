package querypackqueries

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryPackQueriesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewQueryPackQueriesClientWithBaseURI(endpoint string) QueryPackQueriesClient {
	return QueryPackQueriesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

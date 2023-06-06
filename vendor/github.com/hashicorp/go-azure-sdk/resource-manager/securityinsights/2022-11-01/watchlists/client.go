package watchlists

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatchlistsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWatchlistsClientWithBaseURI(endpoint string) WatchlistsClient {
	return WatchlistsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

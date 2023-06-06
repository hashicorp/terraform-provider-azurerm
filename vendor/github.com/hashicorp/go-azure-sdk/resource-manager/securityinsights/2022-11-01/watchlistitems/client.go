package watchlistitems

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatchlistItemsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWatchlistItemsClientWithBaseURI(endpoint string) WatchlistItemsClient {
	return WatchlistItemsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

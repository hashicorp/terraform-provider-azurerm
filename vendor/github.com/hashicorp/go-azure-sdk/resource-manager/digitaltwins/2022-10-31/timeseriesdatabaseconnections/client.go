package timeseriesdatabaseconnections

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimeSeriesDatabaseConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTimeSeriesDatabaseConnectionsClientWithBaseURI(endpoint string) TimeSeriesDatabaseConnectionsClient {
	return TimeSeriesDatabaseConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package eventhubsclusters

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubsClustersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventHubsClustersClientWithBaseURI(endpoint string) EventHubsClustersClient {
	return EventHubsClustersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package eventsources

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourcesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventSourcesClientWithBaseURI(endpoint string) EventSourcesClient {
	return EventSourcesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

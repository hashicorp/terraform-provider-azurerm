package subscriptions

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSubscriptionsClientWithBaseURI(endpoint string) SubscriptionsClient {
	return SubscriptionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

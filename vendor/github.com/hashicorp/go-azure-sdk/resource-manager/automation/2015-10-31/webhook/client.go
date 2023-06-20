package webhook

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWebhookClientWithBaseURI(endpoint string) WebhookClient {
	return WebhookClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

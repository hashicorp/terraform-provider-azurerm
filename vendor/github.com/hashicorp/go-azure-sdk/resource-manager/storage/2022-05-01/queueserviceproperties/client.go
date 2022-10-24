package queueserviceproperties

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueueServicePropertiesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewQueueServicePropertiesClientWithBaseURI(endpoint string) QueueServicePropertiesClient {
	return QueueServicePropertiesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

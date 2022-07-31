package consumergroups

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConsumerGroupsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConsumerGroupsClientWithBaseURI(endpoint string) ConsumerGroupsClient {
	return ConsumerGroupsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

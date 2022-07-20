package topics

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTopicsClientWithBaseURI(endpoint string) TopicsClient {
	return TopicsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

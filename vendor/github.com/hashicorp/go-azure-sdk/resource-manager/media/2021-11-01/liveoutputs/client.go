package liveoutputs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveOutputsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLiveOutputsClientWithBaseURI(endpoint string) LiveOutputsClient {
	return LiveOutputsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

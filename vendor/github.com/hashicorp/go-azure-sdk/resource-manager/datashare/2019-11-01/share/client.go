package share

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ShareClient struct {
	Client  autorest.Client
	baseUri string
}

func NewShareClientWithBaseURI(endpoint string) ShareClient {
	return ShareClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

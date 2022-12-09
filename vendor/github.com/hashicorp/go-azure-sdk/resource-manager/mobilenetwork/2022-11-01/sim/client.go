package sim

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SIMClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSIMClientWithBaseURI(endpoint string) SIMClient {
	return SIMClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

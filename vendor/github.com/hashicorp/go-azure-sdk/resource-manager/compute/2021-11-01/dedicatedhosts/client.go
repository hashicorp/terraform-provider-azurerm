package dedicatedhosts

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDedicatedHostsClientWithBaseURI(endpoint string) DedicatedHostsClient {
	return DedicatedHostsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

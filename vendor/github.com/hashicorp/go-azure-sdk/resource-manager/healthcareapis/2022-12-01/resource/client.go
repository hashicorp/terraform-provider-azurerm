package resource

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewResourceClientWithBaseURI(endpoint string) ResourceClient {
	return ResourceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

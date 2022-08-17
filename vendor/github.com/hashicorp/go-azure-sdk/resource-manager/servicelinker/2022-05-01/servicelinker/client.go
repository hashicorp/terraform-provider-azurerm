package servicelinker

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicelinkerClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServicelinkerClientWithBaseURI(endpoint string) ServicelinkerClient {
	return ServicelinkerClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

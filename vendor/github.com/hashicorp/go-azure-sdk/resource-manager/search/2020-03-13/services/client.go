package services

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServicesClientWithBaseURI(endpoint string) ServicesClient {
	return ServicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

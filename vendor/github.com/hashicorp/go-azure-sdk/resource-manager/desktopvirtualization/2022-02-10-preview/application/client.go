package application

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationClient struct {
	Client  autorest.Client
	baseUri string
}

func NewApplicationClientWithBaseURI(endpoint string) ApplicationClient {
	return ApplicationClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

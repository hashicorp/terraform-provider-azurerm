package linkedservices

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLinkedServicesClientWithBaseURI(endpoint string) LinkedServicesClient {
	return LinkedServicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

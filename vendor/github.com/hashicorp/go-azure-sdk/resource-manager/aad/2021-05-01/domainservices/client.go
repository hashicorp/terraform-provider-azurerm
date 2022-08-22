package domainservices

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainServicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDomainServicesClientWithBaseURI(endpoint string) DomainServicesClient {
	return DomainServicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

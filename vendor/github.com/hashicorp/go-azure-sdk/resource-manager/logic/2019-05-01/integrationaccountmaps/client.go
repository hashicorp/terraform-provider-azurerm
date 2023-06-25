package integrationaccountmaps

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountMapsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIntegrationAccountMapsClientWithBaseURI(endpoint string) IntegrationAccountMapsClient {
	return IntegrationAccountMapsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

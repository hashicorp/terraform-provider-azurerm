package integrationaccountpartners

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountPartnersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIntegrationAccountPartnersClientWithBaseURI(endpoint string) IntegrationAccountPartnersClient {
	return IntegrationAccountPartnersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

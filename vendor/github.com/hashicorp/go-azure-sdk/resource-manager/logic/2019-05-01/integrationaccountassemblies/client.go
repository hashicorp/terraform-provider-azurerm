package integrationaccountassemblies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountAssembliesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIntegrationAccountAssembliesClientWithBaseURI(endpoint string) IntegrationAccountAssembliesClient {
	return IntegrationAccountAssembliesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

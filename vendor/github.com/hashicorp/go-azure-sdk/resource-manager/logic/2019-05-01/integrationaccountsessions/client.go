package integrationaccountsessions

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountSessionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIntegrationAccountSessionsClientWithBaseURI(endpoint string) IntegrationAccountSessionsClient {
	return IntegrationAccountSessionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

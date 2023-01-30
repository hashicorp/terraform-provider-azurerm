package integrationaccountbatchconfigurations

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountBatchConfigurationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIntegrationAccountBatchConfigurationsClientWithBaseURI(endpoint string) IntegrationAccountBatchConfigurationsClient {
	return IntegrationAccountBatchConfigurationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

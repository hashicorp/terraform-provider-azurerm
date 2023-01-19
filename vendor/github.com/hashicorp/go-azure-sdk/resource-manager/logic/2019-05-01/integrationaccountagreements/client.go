package integrationaccountagreements

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountAgreementsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIntegrationAccountAgreementsClientWithBaseURI(endpoint string) IntegrationAccountAgreementsClient {
	return IntegrationAccountAgreementsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

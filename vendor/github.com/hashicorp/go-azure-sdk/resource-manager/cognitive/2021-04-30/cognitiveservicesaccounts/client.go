package cognitiveservicesaccounts

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CognitiveServicesAccountsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCognitiveServicesAccountsClientWithBaseURI(endpoint string) CognitiveServicesAccountsClient {
	return CognitiveServicesAccountsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

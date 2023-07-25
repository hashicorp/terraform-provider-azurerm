package azureadadministrators

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureADAdministratorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAzureADAdministratorsClientWithBaseURI(endpoint string) AzureADAdministratorsClient {
	return AzureADAdministratorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package azuretrafficcollectors

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureTrafficCollectorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAzureTrafficCollectorsClientWithBaseURI(endpoint string) AzureTrafficCollectorsClient {
	return AzureTrafficCollectorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

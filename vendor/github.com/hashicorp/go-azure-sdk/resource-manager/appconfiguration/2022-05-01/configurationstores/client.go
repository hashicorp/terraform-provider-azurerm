package configurationstores

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationStoresClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConfigurationStoresClientWithBaseURI(endpoint string) ConfigurationStoresClient {
	return ConfigurationStoresClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

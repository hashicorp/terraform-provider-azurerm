package configurations

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConfigurationsClientWithBaseURI(endpoint string) ConfigurationsClient {
	return ConfigurationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

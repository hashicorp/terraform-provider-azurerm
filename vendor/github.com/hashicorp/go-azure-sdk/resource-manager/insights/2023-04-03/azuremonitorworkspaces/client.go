package azuremonitorworkspaces

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMonitorWorkspacesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAzureMonitorWorkspacesClientWithBaseURI(endpoint string) AzureMonitorWorkspacesClient {
	return AzureMonitorWorkspacesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

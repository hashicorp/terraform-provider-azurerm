package operationsinacluster

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationsInAClusterClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOperationsInAClusterClientWithBaseURI(endpoint string) OperationsInAClusterClient {
	return OperationsInAClusterClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

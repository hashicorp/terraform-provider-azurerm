package netappaccounts

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetAppAccountsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNetAppAccountsClientWithBaseURI(endpoint string) NetAppAccountsClient {
	return NetAppAccountsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

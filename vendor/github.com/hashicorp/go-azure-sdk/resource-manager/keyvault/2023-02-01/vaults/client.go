package vaults

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVaultsClientWithBaseURI(endpoint string) VaultsClient {
	return VaultsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

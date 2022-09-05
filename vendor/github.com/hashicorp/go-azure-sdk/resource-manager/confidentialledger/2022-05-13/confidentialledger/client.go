package confidentialledger

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfidentialLedgerClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConfidentialLedgerClientWithBaseURI(endpoint string) ConfidentialLedgerClient {
	return ConfidentialLedgerClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

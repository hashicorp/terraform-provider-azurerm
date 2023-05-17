package credential

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCredentialClientWithBaseURI(endpoint string) CredentialClient {
	return CredentialClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package encryptionscopes

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionScopesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEncryptionScopesClientWithBaseURI(endpoint string) EncryptionScopesClient {
	return EncryptionScopesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

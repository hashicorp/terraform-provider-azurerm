package protecteditems

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectedItemsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewProtectedItemsClientWithBaseURI(endpoint string) ProtectedItemsClient {
	return ProtectedItemsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

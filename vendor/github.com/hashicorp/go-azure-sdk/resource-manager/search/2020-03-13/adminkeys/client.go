package adminkeys

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdminKeysClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAdminKeysClientWithBaseURI(endpoint string) AdminKeysClient {
	return AdminKeysClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

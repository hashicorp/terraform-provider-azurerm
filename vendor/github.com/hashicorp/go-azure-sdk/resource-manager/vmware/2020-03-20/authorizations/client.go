package authorizations

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAuthorizationsClientWithBaseURI(endpoint string) AuthorizationsClient {
	return AuthorizationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

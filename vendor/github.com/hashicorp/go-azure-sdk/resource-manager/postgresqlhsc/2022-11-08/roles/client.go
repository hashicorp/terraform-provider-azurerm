package roles

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RolesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRolesClientWithBaseURI(endpoint string) RolesClient {
	return RolesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

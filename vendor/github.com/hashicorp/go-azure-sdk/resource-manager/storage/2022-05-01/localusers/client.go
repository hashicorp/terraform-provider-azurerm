package localusers

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalUsersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLocalUsersClientWithBaseURI(endpoint string) LocalUsersClient {
	return LocalUsersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

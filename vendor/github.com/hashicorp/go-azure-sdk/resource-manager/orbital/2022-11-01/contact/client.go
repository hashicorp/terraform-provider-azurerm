package contact

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactClient struct {
	Client  autorest.Client
	baseUri string
}

func NewContactClientWithBaseURI(endpoint string) ContactClient {
	return ContactClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

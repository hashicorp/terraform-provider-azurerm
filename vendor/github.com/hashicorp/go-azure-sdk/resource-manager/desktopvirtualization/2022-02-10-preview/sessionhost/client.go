package sessionhost

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SessionHostClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSessionHostClientWithBaseURI(endpoint string) SessionHostClient {
	return SessionHostClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

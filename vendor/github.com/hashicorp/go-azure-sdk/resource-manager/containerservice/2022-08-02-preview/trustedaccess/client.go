package trustedaccess

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrustedAccessClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTrustedAccessClientWithBaseURI(endpoint string) TrustedAccessClient {
	return TrustedAccessClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package simgroup

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SIMGroupClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSIMGroupClientWithBaseURI(endpoint string) SIMGroupClient {
	return SIMGroupClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

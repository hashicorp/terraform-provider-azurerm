package encodings

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncodingsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEncodingsClientWithBaseURI(endpoint string) EncodingsClient {
	return EncodingsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

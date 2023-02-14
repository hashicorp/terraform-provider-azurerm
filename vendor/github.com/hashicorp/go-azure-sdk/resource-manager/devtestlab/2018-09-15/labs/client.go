package labs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLabsClientWithBaseURI(endpoint string) LabsClient {
	return LabsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package batchaccount

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchAccountClient struct {
	Client  autorest.Client
	baseUri string
}

func NewBatchAccountClientWithBaseURI(endpoint string) BatchAccountClient {
	return BatchAccountClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

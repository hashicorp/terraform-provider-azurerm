package functions

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFunctionsClientWithBaseURI(endpoint string) FunctionsClient {
	return FunctionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

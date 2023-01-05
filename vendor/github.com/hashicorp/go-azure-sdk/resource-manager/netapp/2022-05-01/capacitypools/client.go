package capacitypools

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityPoolsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCapacityPoolsClientWithBaseURI(endpoint string) CapacityPoolsClient {
	return CapacityPoolsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

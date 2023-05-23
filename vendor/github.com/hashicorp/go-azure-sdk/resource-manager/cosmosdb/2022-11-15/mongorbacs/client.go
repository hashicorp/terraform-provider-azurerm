package mongorbacs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongorbacsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMongorbacsClientWithBaseURI(endpoint string) MongorbacsClient {
	return MongorbacsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

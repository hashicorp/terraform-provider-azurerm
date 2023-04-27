package iotconnectors

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotConnectorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIotConnectorsClientWithBaseURI(endpoint string) IotConnectorsClient {
	return IotConnectorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

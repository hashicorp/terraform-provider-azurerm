package amlfilesystems

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAmlFilesystemsClientWithBaseURI(endpoint string) AmlFilesystemsClient {
	return AmlFilesystemsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

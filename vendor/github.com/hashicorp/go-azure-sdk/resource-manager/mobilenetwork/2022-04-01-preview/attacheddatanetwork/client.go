package attacheddatanetwork

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttachedDataNetworkClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAttachedDataNetworkClientWithBaseURI(endpoint string) AttachedDataNetworkClient {
	return AttachedDataNetworkClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

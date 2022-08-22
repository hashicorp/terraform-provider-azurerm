package dedicatedhostgroups

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostGroupsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDedicatedHostGroupsClientWithBaseURI(endpoint string) DedicatedHostGroupsClient {
	return DedicatedHostGroupsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

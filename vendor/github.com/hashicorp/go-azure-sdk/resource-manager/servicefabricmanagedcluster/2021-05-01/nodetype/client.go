package nodetype

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNodeTypeClientWithBaseURI(endpoint string) NodeTypeClient {
	return NodeTypeClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

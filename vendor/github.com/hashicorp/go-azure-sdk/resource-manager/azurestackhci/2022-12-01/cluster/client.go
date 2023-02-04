package cluster

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterClient struct {
	Client  autorest.Client
	baseUri string
}

func NewClusterClientWithBaseURI(endpoint string) ClusterClient {
	return ClusterClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

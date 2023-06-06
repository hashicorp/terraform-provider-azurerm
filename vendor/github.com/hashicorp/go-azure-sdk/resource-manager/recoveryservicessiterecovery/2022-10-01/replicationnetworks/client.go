package replicationnetworks

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationNetworksClient struct {
	Client  autorest.Client
	baseUri string
}

func NewReplicationNetworksClientWithBaseURI(endpoint string) ReplicationNetworksClient {
	return ReplicationNetworksClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package replicationprotecteditems

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProtectedItemsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewReplicationProtectedItemsClientWithBaseURI(endpoint string) ReplicationProtectedItemsClient {
	return ReplicationProtectedItemsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

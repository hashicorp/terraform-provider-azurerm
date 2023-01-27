package managedclustersnapshots

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSnapshotsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedClusterSnapshotsClientWithBaseURI(endpoint string) ManagedClusterSnapshotsClient {
	return ManagedClusterSnapshotsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package snapshotpolicy

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotPolicyClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSnapshotPolicyClientWithBaseURI(endpoint string) SnapshotPolicyClient {
	return SnapshotPolicyClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

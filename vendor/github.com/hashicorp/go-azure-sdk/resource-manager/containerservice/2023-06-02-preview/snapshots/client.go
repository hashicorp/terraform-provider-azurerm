package snapshots

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSnapshotsClientWithBaseURI(endpoint string) SnapshotsClient {
	return SnapshotsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

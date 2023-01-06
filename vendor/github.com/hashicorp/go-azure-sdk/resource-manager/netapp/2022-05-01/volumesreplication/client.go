package volumesreplication

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumesReplicationClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVolumesReplicationClientWithBaseURI(endpoint string) VolumesReplicationClient {
	return VolumesReplicationClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

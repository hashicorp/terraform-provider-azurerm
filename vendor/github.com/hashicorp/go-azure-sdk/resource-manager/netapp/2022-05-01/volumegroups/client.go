package volumegroups

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeGroupsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVolumeGroupsClientWithBaseURI(endpoint string) VolumeGroupsClient {
	return VolumeGroupsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

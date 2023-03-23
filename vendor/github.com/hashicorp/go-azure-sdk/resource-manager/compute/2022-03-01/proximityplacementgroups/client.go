package proximityplacementgroups

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProximityPlacementGroupsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewProximityPlacementGroupsClientWithBaseURI(endpoint string) ProximityPlacementGroupsClient {
	return ProximityPlacementGroupsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

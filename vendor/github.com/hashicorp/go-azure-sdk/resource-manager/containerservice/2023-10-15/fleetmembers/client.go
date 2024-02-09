package fleetmembers

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetMembersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFleetMembersClientWithBaseURI(endpoint string) FleetMembersClient {
	return FleetMembersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

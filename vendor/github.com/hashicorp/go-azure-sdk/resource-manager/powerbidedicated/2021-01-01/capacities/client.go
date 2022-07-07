package capacities

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacitiesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCapacitiesClientWithBaseURI(endpoint string) CapacitiesClient {
	return CapacitiesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

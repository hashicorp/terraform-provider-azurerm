package fleetupdatestrategies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetUpdateStrategiesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFleetUpdateStrategiesClientWithBaseURI(endpoint string) FleetUpdateStrategiesClient {
	return FleetUpdateStrategiesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

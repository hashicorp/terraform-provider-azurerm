package fluidrelaycontainers

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluidRelayContainersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFluidRelayContainersClientWithBaseURI(endpoint string) FluidRelayContainersClient {
	return FluidRelayContainersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

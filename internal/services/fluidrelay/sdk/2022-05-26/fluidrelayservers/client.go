package fluidrelayservers

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluidRelayServersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFluidRelayServersClientWithBaseURI(endpoint string) FluidRelayServersClient {
	return FluidRelayServersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package outboundnetworkdependenciesendpoints

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutboundNetworkDependenciesEndpointsClient struct {
	Client *resourcemanager.Client
}

func NewOutboundNetworkDependenciesEndpointsClientWithBaseURI(api environments.Api) (*OutboundNetworkDependenciesEndpointsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "outboundnetworkdependenciesendpoints", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OutboundNetworkDependenciesEndpointsClient: %+v", err)
	}

	return &OutboundNetworkDependenciesEndpointsClient{
		Client: client,
	}, nil
}

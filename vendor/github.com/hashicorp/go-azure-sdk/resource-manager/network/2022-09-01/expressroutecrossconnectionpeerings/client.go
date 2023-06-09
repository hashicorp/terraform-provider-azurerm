package expressroutecrossconnectionpeerings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCrossConnectionPeeringsClient struct {
	Client *resourcemanager.Client
}

func NewExpressRouteCrossConnectionPeeringsClientWithBaseURI(api environments.Api) (*ExpressRouteCrossConnectionPeeringsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "expressroutecrossconnectionpeerings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRouteCrossConnectionPeeringsClient: %+v", err)
	}

	return &ExpressRouteCrossConnectionPeeringsClient{
		Client: client,
	}, nil
}

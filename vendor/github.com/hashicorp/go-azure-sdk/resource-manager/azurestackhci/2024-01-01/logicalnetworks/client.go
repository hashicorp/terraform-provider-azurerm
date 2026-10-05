package logicalnetworks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogicalNetworksClient struct {
	Client *resourcemanager.Client
}

func NewLogicalNetworksClientWithBaseURI(sdkApi sdkEnv.Api) (*LogicalNetworksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "logicalnetworks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LogicalNetworksClient: %+v", err)
	}

	return &LogicalNetworksClient{
		Client: client,
	}, nil
}

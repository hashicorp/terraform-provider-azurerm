package poolchange

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolChangeClient struct {
	Client *resourcemanager.Client
}

func NewPoolChangeClientWithBaseURI(sdkApi sdkEnv.Api) (*PoolChangeClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "poolchange", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PoolChangeClient: %+v", err)
	}

	return &PoolChangeClient{
		Client: client,
	}, nil
}

package pools

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolsClient struct {
	Client *resourcemanager.Client
}

func NewPoolsClientWithBaseURI(sdkApi sdkEnv.Api) (*PoolsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "pools", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PoolsClient: %+v", err)
	}

	return &PoolsClient{
		Client: client,
	}, nil
}

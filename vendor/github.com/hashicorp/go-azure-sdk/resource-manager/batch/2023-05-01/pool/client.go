package pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolClient struct {
	Client *resourcemanager.Client
}

func NewPoolClientWithBaseURI(sdkApi sdkEnv.Api) (*PoolClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "pool", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PoolClient: %+v", err)
	}

	return &PoolClient{
		Client: client,
	}, nil
}

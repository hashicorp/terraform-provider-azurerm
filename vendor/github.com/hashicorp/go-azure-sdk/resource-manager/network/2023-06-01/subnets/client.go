package subnets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubnetsClient struct {
	Client *resourcemanager.Client
}

func NewSubnetsClientWithBaseURI(sdkApi sdkEnv.Api) (*SubnetsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "subnets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SubnetsClient: %+v", err)
	}

	return &SubnetsClient{
		Client: client,
	}, nil
}

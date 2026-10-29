package nodepools

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodePoolsClient struct {
	Client *resourcemanager.Client
}

func NewNodePoolsClientWithBaseURI(sdkApi sdkEnv.Api) (*NodePoolsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nodepools", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NodePoolsClient: %+v", err)
	}

	return &NodePoolsClient{
		Client: client,
	}, nil
}

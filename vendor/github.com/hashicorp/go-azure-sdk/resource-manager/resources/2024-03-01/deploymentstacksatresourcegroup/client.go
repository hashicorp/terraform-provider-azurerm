package deploymentstacksatresourcegroup

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentStacksAtResourceGroupClient struct {
	Client *resourcemanager.Client
}

func NewDeploymentStacksAtResourceGroupClientWithBaseURI(sdkApi sdkEnv.Api) (*DeploymentStacksAtResourceGroupClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "deploymentstacksatresourcegroup", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeploymentStacksAtResourceGroupClient: %+v", err)
	}

	return &DeploymentStacksAtResourceGroupClient{
		Client: client,
	}, nil
}

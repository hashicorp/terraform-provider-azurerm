package deploymentstacksatmanagementgroup

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentStacksAtManagementGroupClient struct {
	Client *resourcemanager.Client
}

func NewDeploymentStacksAtManagementGroupClientWithBaseURI(sdkApi sdkEnv.Api) (*DeploymentStacksAtManagementGroupClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "deploymentstacksatmanagementgroup", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeploymentStacksAtManagementGroupClient: %+v", err)
	}

	return &DeploymentStacksAtManagementGroupClient{
		Client: client,
	}, nil
}

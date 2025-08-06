package deploymentscripts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentScriptsClient struct {
	Client *resourcemanager.Client
}

func NewDeploymentScriptsClientWithBaseURI(sdkApi sdkEnv.Api) (*DeploymentScriptsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "deploymentscripts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeploymentScriptsClient: %+v", err)
	}

	return &DeploymentScriptsClient{
		Client: client,
	}, nil
}

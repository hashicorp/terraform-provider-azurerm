package deploymentsafeguards

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentSafeguardsClient struct {
	Client *resourcemanager.Client
}

func NewDeploymentSafeguardsClientWithBaseURI(sdkApi sdkEnv.Api) (*DeploymentSafeguardsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "deploymentsafeguards", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeploymentSafeguardsClient: %+v", err)
	}

	return &DeploymentSafeguardsClient{
		Client: client,
	}, nil
}

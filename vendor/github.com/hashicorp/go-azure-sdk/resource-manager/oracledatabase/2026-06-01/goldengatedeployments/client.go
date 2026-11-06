package goldengatedeployments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GoldenGateDeploymentsClient struct {
	Client *resourcemanager.Client
}

func NewGoldenGateDeploymentsClientWithBaseURI(sdkApi sdkEnv.Api) (*GoldenGateDeploymentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "goldengatedeployments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GoldenGateDeploymentsClient: %+v", err)
	}

	return &GoldenGateDeploymentsClient{
		Client: client,
	}, nil
}

package integrationruntimes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimesClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationRuntimesClientWithBaseURI(sdkApi sdkEnv.Api) (*IntegrationRuntimesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "integrationruntimes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationRuntimesClient: %+v", err)
	}

	return &IntegrationRuntimesClient{
		Client: client,
	}, nil
}

package integrationruntimeenableinteractivequery

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeEnableInteractiveQueryClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationRuntimeEnableInteractiveQueryClientWithBaseURI(sdkApi sdkEnv.Api) (*IntegrationRuntimeEnableInteractiveQueryClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "integrationruntimeenableinteractivequery", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationRuntimeEnableInteractiveQueryClient: %+v", err)
	}

	return &IntegrationRuntimeEnableInteractiveQueryClient{
		Client: client,
	}, nil
}

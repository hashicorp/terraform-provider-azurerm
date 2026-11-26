package integrationruntimedisableinteractivequery

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeDisableInteractiveQueryClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationRuntimeDisableInteractiveQueryClientWithBaseURI(sdkApi sdkEnv.Api) (*IntegrationRuntimeDisableInteractiveQueryClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "integrationruntimedisableinteractivequery", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationRuntimeDisableInteractiveQueryClient: %+v", err)
	}

	return &IntegrationRuntimeDisableInteractiveQueryClient{
		Client: client,
	}, nil
}

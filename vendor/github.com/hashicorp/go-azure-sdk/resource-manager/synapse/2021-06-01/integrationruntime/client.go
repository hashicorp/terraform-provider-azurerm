package integrationruntime

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationRuntimeClientWithBaseURI(sdkApi sdkEnv.Api) (*IntegrationRuntimeClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "integrationruntime", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationRuntimeClient: %+v", err)
	}

	return &IntegrationRuntimeClient{
		Client: client,
	}, nil
}

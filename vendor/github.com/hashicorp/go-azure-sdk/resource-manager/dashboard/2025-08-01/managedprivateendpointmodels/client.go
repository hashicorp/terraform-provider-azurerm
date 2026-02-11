package managedprivateendpointmodels

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedPrivateEndpointModelsClient struct {
	Client *resourcemanager.Client
}

func NewManagedPrivateEndpointModelsClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedPrivateEndpointModelsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedprivateendpointmodels", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedPrivateEndpointModelsClient: %+v", err)
	}

	return &ManagedPrivateEndpointModelsClient{
		Client: client,
	}, nil
}

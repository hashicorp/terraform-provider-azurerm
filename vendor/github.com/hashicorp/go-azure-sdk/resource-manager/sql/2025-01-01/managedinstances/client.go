package managedinstances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstancesClient struct {
	Client *resourcemanager.Client
}

func NewManagedInstancesClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedInstancesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedinstances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedInstancesClient: %+v", err)
	}

	return &ManagedInstancesClient{
		Client: client,
	}, nil
}

package provisionedclusterinstances

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisionedClusterInstancesClient struct {
	Client *resourcemanager.Client
}

func NewProvisionedClusterInstancesClientWithBaseURI(sdkApi sdkEnv.Api) (*ProvisionedClusterInstancesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "provisionedclusterinstances", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProvisionedClusterInstancesClient: %+v", err)
	}

	return &ProvisionedClusterInstancesClient{
		Client: client,
	}, nil
}

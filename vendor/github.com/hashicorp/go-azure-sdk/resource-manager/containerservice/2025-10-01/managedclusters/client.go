package managedclusters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClustersClient struct {
	Client *resourcemanager.Client
}

func NewManagedClustersClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedClustersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedclusters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedClustersClient: %+v", err)
	}

	return &ManagedClustersClient{
		Client: client,
	}, nil
}

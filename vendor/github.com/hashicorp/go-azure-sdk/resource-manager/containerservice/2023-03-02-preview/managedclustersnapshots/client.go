package managedclustersnapshots

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSnapshotsClient struct {
	Client *resourcemanager.Client
}

func NewManagedClusterSnapshotsClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedClusterSnapshotsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedclustersnapshots", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedClusterSnapshotsClient: %+v", err)
	}

	return &ManagedClusterSnapshotsClient{
		Client: client,
	}, nil
}

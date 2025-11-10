package managedcluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterClient struct {
	Client *resourcemanager.Client
}

func NewManagedClusterClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedClusterClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedcluster", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedClusterClient: %+v", err)
	}

	return &ManagedClusterClient{
		Client: client,
	}, nil
}

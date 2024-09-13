package replicationnetworkmappings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationNetworkMappingsClient struct {
	Client *resourcemanager.Client
}

func NewReplicationNetworkMappingsClientWithBaseURI(sdkApi sdkEnv.Api) (*ReplicationNetworkMappingsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "replicationnetworkmappings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReplicationNetworkMappingsClient: %+v", err)
	}

	return &ReplicationNetworkMappingsClient{
		Client: client,
	}, nil
}

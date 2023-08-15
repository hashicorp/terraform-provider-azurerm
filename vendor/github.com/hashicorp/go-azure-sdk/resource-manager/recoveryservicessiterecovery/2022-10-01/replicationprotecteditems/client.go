package replicationprotecteditems

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProtectedItemsClient struct {
	Client *resourcemanager.Client
}

func NewReplicationProtectedItemsClientWithBaseURI(sdkApi sdkEnv.Api) (*ReplicationProtectedItemsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "replicationprotecteditems", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReplicationProtectedItemsClient: %+v", err)
	}

	return &ReplicationProtectedItemsClient{
		Client: client,
	}, nil
}

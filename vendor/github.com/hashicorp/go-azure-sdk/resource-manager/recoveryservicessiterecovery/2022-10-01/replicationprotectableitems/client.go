package replicationprotectableitems

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProtectableItemsClient struct {
	Client *resourcemanager.Client
}

func NewReplicationProtectableItemsClientWithBaseURI(sdkApi sdkEnv.Api) (*ReplicationProtectableItemsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "replicationprotectableitems", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReplicationProtectableItemsClient: %+v", err)
	}

	return &ReplicationProtectableItemsClient{
		Client: client,
	}, nil
}

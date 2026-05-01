package storagequeues

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageQueuesClient struct {
	Client *resourcemanager.Client
}

func NewStorageQueuesClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageQueuesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "storagequeues", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageQueuesClient: %+v", err)
	}

	return &StorageQueuesClient{
		Client: client,
	}, nil
}

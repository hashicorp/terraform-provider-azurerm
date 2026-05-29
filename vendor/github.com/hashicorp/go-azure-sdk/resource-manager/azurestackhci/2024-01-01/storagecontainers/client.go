package storagecontainers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageContainersClient struct {
	Client *resourcemanager.Client
}

func NewStorageContainersClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageContainersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "storagecontainers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageContainersClient: %+v", err)
	}

	return &StorageContainersClient{
		Client: client,
	}, nil
}

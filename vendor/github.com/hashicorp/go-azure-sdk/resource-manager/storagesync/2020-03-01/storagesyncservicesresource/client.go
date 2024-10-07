package storagesyncservicesresource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageSyncServicesResourceClient struct {
	Client *resourcemanager.Client
}

func NewStorageSyncServicesResourceClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageSyncServicesResourceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "storagesyncservicesresource", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageSyncServicesResourceClient: %+v", err)
	}

	return &StorageSyncServicesResourceClient{
		Client: client,
	}, nil
}

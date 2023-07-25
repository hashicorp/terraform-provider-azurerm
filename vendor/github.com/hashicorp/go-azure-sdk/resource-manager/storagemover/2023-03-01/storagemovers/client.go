package storagemovers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageMoversClient struct {
	Client *resourcemanager.Client
}

func NewStorageMoversClientWithBaseURI(api environments.Api) (*StorageMoversClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "storagemovers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageMoversClient: %+v", err)
	}

	return &StorageMoversClient{
		Client: client,
	}, nil
}

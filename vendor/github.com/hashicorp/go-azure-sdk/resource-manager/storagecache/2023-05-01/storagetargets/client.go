package storagetargets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTargetsClient struct {
	Client *resourcemanager.Client
}

func NewStorageTargetsClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageTargetsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "storagetargets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageTargetsClient: %+v", err)
	}

	return &StorageTargetsClient{
		Client: client,
	}, nil
}

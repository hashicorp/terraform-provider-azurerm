package storagemovers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageMoversClient struct {
	Client *resourcemanager.Client
}

func NewStorageMoversClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageMoversClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "storagemovers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageMoversClient: %+v", err)
	}

	return &StorageMoversClient{
		Client: client,
	}, nil
}

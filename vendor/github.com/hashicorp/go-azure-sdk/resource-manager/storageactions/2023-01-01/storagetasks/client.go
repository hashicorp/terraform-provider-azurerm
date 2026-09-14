package storagetasks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTasksClient struct {
	Client *resourcemanager.Client
}

func NewStorageTasksClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageTasksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "storagetasks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageTasksClient: %+v", err)
	}

	return &StorageTasksClient{
		Client: client,
	}, nil
}

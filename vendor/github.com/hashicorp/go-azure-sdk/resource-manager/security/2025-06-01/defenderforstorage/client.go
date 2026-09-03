package defenderforstorage

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefenderForStorageClient struct {
	Client *resourcemanager.Client
}

func NewDefenderForStorageClientWithBaseURI(sdkApi sdkEnv.Api) (*DefenderForStorageClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "defenderforstorage", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DefenderForStorageClient: %+v", err)
	}

	return &DefenderForStorageClient{
		Client: client,
	}, nil
}

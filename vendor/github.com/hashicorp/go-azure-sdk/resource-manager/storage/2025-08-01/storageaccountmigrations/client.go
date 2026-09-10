package storageaccountmigrations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountMigrationsClient struct {
	Client *resourcemanager.Client
}

func NewStorageAccountMigrationsClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageAccountMigrationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "storageaccountmigrations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageAccountMigrationsClient: %+v", err)
	}

	return &StorageAccountMigrationsClient{
		Client: client,
	}, nil
}

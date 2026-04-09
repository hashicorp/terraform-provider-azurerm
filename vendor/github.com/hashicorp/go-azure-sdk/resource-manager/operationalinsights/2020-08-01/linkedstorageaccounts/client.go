package linkedstorageaccounts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedStorageAccountsClient struct {
	Client *resourcemanager.Client
}

func NewLinkedStorageAccountsClientWithBaseURI(sdkApi sdkEnv.Api) (*LinkedStorageAccountsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "linkedstorageaccounts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LinkedStorageAccountsClient: %+v", err)
	}

	return &LinkedStorageAccountsClient{
		Client: client,
	}, nil
}

package inventoryitems

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InventoryItemsClient struct {
	Client *resourcemanager.Client
}

func NewInventoryItemsClientWithBaseURI(sdkApi sdkEnv.Api) (*InventoryItemsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "inventoryitems", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating InventoryItemsClient: %+v", err)
	}

	return &InventoryItemsClient{
		Client: client,
	}, nil
}

package watchlistitems

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatchlistItemsClient struct {
	Client *resourcemanager.Client
}

func NewWatchlistItemsClientWithBaseURI(sdkApi sdkEnv.Api) (*WatchlistItemsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "watchlistitems", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WatchlistItemsClient: %+v", err)
	}

	return &WatchlistItemsClient{
		Client: client,
	}, nil
}

package watchlists

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatchlistsClient struct {
	Client *resourcemanager.Client
}

func NewWatchlistsClientWithBaseURI(sdkApi sdkEnv.Api) (*WatchlistsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "watchlists", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WatchlistsClient: %+v", err)
	}

	return &WatchlistsClient{
		Client: client,
	}, nil
}

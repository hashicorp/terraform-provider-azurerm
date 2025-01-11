package watcher

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatcherClient struct {
	Client *resourcemanager.Client
}

func NewWatcherClientWithBaseURI(sdkApi sdkEnv.Api) (*WatcherClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "watcher", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating WatcherClient: %+v", err)
	}

	return &WatcherClient{
		Client: client,
	}, nil
}

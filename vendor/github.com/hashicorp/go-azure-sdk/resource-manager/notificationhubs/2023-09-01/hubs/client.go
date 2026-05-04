package hubs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubsClient struct {
	Client *resourcemanager.Client
}

func NewHubsClientWithBaseURI(sdkApi sdkEnv.Api) (*HubsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "hubs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating HubsClient: %+v", err)
	}

	return &HubsClient{
		Client: client,
	}, nil
}

package eventhubs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubsClient struct {
	Client *resourcemanager.Client
}

func NewEventHubsClientWithBaseURI(sdkApi sdkEnv.Api) (*EventHubsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "eventhubs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EventHubsClient: %+v", err)
	}

	return &EventHubsClient{
		Client: client,
	}, nil
}

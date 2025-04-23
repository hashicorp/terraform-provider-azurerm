package eventhubsclusters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubsClustersClient struct {
	Client *resourcemanager.Client
}

func NewEventHubsClustersClientWithBaseURI(sdkApi sdkEnv.Api) (*EventHubsClustersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "eventhubsclusters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EventHubsClustersClient: %+v", err)
	}

	return &EventHubsClustersClient{
		Client: client,
	}, nil
}

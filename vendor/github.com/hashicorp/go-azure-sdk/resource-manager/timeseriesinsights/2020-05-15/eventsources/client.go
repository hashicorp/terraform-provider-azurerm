package eventsources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourcesClient struct {
	Client *resourcemanager.Client
}

func NewEventSourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*EventSourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "eventsources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EventSourcesClient: %+v", err)
	}

	return &EventSourcesClient{
		Client: client,
	}, nil
}

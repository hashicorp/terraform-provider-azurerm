package liveevents

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventsClient struct {
	Client *resourcemanager.Client
}

func NewLiveEventsClientWithBaseURI(sdkApi sdkEnv.Api) (*LiveEventsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "liveevents", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LiveEventsClient: %+v", err)
	}

	return &LiveEventsClient{
		Client: client,
	}, nil
}

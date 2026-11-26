package trafficmanagers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficmanagersClient struct {
	Client *resourcemanager.Client
}

func NewTrafficmanagersClientWithBaseURI(sdkApi sdkEnv.Api) (*TrafficmanagersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "trafficmanagers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TrafficmanagersClient: %+v", err)
	}

	return &TrafficmanagersClient{
		Client: client,
	}, nil
}

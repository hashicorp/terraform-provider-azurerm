package trafficanalytics

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficAnalyticsClient struct {
	Client *resourcemanager.Client
}

func NewTrafficAnalyticsClientWithBaseURI(sdkApi sdkEnv.Api) (*TrafficAnalyticsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "trafficanalytics", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TrafficAnalyticsClient: %+v", err)
	}

	return &TrafficAnalyticsClient{
		Client: client,
	}, nil
}

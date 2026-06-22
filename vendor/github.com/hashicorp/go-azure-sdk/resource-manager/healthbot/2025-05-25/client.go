package v2025_05_25

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25/healthbots"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	HealthBots *healthbots.HealthBotsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	healthBotsClient, err := healthbots.NewHealthBotsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building HealthBots client: %+v", err)
	}
	configureFunc(healthBotsClient.Client)

	return &Client{
		HealthBots: healthBotsClient,
	}, nil
}

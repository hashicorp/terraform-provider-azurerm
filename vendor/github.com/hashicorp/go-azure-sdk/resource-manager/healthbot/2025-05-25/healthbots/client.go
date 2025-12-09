package healthbots

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthBotsClient struct {
	Client *resourcemanager.Client
}

func NewHealthBotsClientWithBaseURI(sdkApi sdkEnv.Api) (*HealthBotsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "healthbots", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating HealthBotsClient: %+v", err)
	}

	return &HealthBotsClient{
		Client: client,
	}, nil
}

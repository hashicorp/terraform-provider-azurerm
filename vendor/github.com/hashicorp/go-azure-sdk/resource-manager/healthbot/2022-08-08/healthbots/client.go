package healthbots

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthbotsClient struct {
	Client *resourcemanager.Client
}

func NewHealthbotsClientWithBaseURI(sdkApi sdkEnv.Api) (*HealthbotsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "healthbots", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating HealthbotsClient: %+v", err)
	}

	return &HealthbotsClient{
		Client: client,
	}, nil
}

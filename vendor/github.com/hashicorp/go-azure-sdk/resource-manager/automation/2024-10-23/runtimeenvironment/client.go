package runtimeenvironment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuntimeEnvironmentClient struct {
	Client *resourcemanager.Client
}

func NewRuntimeEnvironmentClientWithBaseURI(sdkApi sdkEnv.Api) (*RuntimeEnvironmentClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "runtimeenvironment", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RuntimeEnvironmentClient: %+v", err)
	}

	return &RuntimeEnvironmentClient{
		Client: client,
	}, nil
}

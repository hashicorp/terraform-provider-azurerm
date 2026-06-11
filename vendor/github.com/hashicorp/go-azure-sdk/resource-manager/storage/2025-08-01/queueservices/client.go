package queueservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueueServicesClient struct {
	Client *resourcemanager.Client
}

func NewQueueServicesClientWithBaseURI(sdkApi sdkEnv.Api) (*QueueServicesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "queueservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating QueueServicesClient: %+v", err)
	}

	return &QueueServicesClient{
		Client: client,
	}, nil
}

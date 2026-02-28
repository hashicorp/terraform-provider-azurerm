package broker

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerClient struct {
	Client *resourcemanager.Client
}

func NewBrokerClientWithBaseURI(sdkApi sdkEnv.Api) (*BrokerClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "broker", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BrokerClient: %+v", err)
	}

	return &BrokerClient{
		Client: client,
	}, nil
}

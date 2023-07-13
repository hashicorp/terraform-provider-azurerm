package iotconnectors

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotConnectorsClient struct {
	Client *resourcemanager.Client
}

func NewIotConnectorsClientWithBaseURI(api environments.Api) (*IotConnectorsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "iotconnectors", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IotConnectorsClient: %+v", err)
	}

	return &IotConnectorsClient{
		Client: client,
	}, nil
}

package trafficcontrollerinterface

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficControllerInterfaceClient struct {
	Client *resourcemanager.Client
}

func NewTrafficControllerInterfaceClientWithBaseURI(api environments.Api) (*TrafficControllerInterfaceClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "trafficcontrollerinterface", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TrafficControllerInterfaceClient: %+v", err)
	}

	return &TrafficControllerInterfaceClient{
		Client: client,
	}, nil
}

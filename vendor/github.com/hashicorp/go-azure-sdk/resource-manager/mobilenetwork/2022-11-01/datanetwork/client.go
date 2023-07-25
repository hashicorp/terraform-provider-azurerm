package datanetwork

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataNetworkClient struct {
	Client *resourcemanager.Client
}

func NewDataNetworkClientWithBaseURI(api environments.Api) (*DataNetworkClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "datanetwork", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataNetworkClient: %+v", err)
	}

	return &DataNetworkClient{
		Client: client,
	}, nil
}

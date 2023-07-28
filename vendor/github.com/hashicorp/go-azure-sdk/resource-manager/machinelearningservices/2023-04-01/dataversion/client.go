package dataversion

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataVersionClient struct {
	Client *resourcemanager.Client
}

func NewDataVersionClientWithBaseURI(api environments.Api) (*DataVersionClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "dataversion", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataVersionClient: %+v", err)
	}

	return &DataVersionClient{
		Client: client,
	}, nil
}

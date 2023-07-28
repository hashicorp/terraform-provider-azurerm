package datacontainer

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataContainerClient struct {
	Client *resourcemanager.Client
}

func NewDataContainerClientWithBaseURI(api environments.Api) (*DataContainerClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "datacontainer", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataContainerClient: %+v", err)
	}

	return &DataContainerClient{
		Client: client,
	}, nil
}

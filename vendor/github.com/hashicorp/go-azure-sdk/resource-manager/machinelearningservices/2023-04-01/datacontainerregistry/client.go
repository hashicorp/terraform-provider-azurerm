package datacontainerregistry

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataContainerRegistryClient struct {
	Client *resourcemanager.Client
}

func NewDataContainerRegistryClientWithBaseURI(api environments.Api) (*DataContainerRegistryClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "datacontainerregistry", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataContainerRegistryClient: %+v", err)
	}

	return &DataContainerRegistryClient{
		Client: client,
	}, nil
}

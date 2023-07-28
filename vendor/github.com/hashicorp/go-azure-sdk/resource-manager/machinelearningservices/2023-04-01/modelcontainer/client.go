package modelcontainer

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModelContainerClient struct {
	Client *resourcemanager.Client
}

func NewModelContainerClientWithBaseURI(api environments.Api) (*ModelContainerClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "modelcontainer", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ModelContainerClient: %+v", err)
	}

	return &ModelContainerClient{
		Client: client,
	}, nil
}

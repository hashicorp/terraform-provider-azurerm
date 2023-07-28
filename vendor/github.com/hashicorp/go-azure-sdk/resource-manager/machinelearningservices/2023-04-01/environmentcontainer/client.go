package environmentcontainer

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentContainerClient struct {
	Client *resourcemanager.Client
}

func NewEnvironmentContainerClientWithBaseURI(api environments.Api) (*EnvironmentContainerClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "environmentcontainer", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EnvironmentContainerClient: %+v", err)
	}

	return &EnvironmentContainerClient{
		Client: client,
	}, nil
}

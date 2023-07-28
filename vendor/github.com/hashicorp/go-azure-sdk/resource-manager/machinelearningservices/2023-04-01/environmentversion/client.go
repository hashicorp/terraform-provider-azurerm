package environmentversion

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentVersionClient struct {
	Client *resourcemanager.Client
}

func NewEnvironmentVersionClientWithBaseURI(api environments.Api) (*EnvironmentVersionClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "environmentversion", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EnvironmentVersionClient: %+v", err)
	}

	return &EnvironmentVersionClient{
		Client: client,
	}, nil
}

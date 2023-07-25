package module

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModuleClient struct {
	Client *resourcemanager.Client
}

func NewModuleClientWithBaseURI(api environments.Api) (*ModuleClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "module", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ModuleClient: %+v", err)
	}

	return &ModuleClient{
		Client: client,
	}, nil
}

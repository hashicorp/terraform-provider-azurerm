package applicationdefinitions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationDefinitionsClient struct {
	Client *resourcemanager.Client
}

func NewApplicationDefinitionsClientWithBaseURI(api environments.Api) (*ApplicationDefinitionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "applicationdefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApplicationDefinitionsClient: %+v", err)
	}

	return &ApplicationDefinitionsClient{
		Client: client,
	}, nil
}

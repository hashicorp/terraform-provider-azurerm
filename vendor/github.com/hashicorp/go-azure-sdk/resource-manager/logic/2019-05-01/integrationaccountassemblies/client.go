package integrationaccountassemblies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountAssembliesClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationAccountAssembliesClientWithBaseURI(api environments.Api) (*IntegrationAccountAssembliesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "integrationaccountassemblies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationAccountAssembliesClient: %+v", err)
	}

	return &IntegrationAccountAssembliesClient{
		Client: client,
	}, nil
}

package availableservicealiases

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableServiceAliasesClient struct {
	Client *resourcemanager.Client
}

func NewAvailableServiceAliasesClientWithBaseURI(api environments.Api) (*AvailableServiceAliasesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "availableservicealiases", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AvailableServiceAliasesClient: %+v", err)
	}

	return &AvailableServiceAliasesClient{
		Client: client,
	}, nil
}

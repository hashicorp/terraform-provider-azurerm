package configurationstores

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationStoresClient struct {
	Client *resourcemanager.Client
}

func NewConfigurationStoresClientWithBaseURI(api environments.Api) (*ConfigurationStoresClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "configurationstores", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConfigurationStoresClient: %+v", err)
	}

	return &ConfigurationStoresClient{
		Client: client,
	}, nil
}

package deletedconfigurationstores

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedConfigurationStoresClient struct {
	Client *resourcemanager.Client
}

func NewDeletedConfigurationStoresClientWithBaseURI(api environments.Api) (*DeletedConfigurationStoresClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "deletedconfigurationstores", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedConfigurationStoresClient: %+v", err)
	}

	return &DeletedConfigurationStoresClient{
		Client: client,
	}, nil
}

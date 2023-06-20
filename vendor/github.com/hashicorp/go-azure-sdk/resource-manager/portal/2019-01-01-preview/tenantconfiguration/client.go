package tenantconfiguration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantConfigurationClient struct {
	Client *resourcemanager.Client
}

func NewTenantConfigurationClientWithBaseURI(api environments.Api) (*TenantConfigurationClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "tenantconfiguration", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TenantConfigurationClient: %+v", err)
	}

	return &TenantConfigurationClient{
		Client: client,
	}, nil
}

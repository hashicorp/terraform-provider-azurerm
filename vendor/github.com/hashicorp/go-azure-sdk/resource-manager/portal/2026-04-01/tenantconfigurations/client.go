package tenantconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewTenantConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*TenantConfigurationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "tenantconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TenantConfigurationsClient: %+v", err)
	}

	return &TenantConfigurationsClient{
		Client: client,
	}, nil
}

package tenantaccess

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantAccessClient struct {
	Client *resourcemanager.Client
}

func NewTenantAccessClientWithBaseURI(sdkApi sdkEnv.Api) (*TenantAccessClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "tenantaccess", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TenantAccessClient: %+v", err)
	}

	return &TenantAccessClient{
		Client: client,
	}, nil
}

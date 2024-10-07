package v2021_04_01_preview

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Tenants *tenants.TenantsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	tenantsClient, err := tenants.NewTenantsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Tenants client: %+v", err)
	}
	configureFunc(tenantsClient.Client)

	return &Client{
		Tenants: tenantsClient,
	}, nil
}

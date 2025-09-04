package v2024_11_30

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/managedidentities"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	ManagedIdentities *managedidentities.ManagedIdentitiesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	managedIdentitiesClient, err := managedidentities.NewManagedIdentitiesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ManagedIdentities client: %+v", err)
	}
	configureFunc(managedIdentitiesClient.Client)

	return &Client{
		ManagedIdentities: managedIdentitiesClient,
	}, nil
}

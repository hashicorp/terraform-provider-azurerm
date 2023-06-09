package v2023_01_31

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2023-01-31/managedidentities"
)

type Client struct {
	ManagedIdentities *managedidentities.ManagedIdentitiesClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	managedIdentitiesClient := managedidentities.NewManagedIdentitiesClientWithBaseURI(endpoint)
	configureAuthFunc(&managedIdentitiesClient.Client)

	return Client{
		ManagedIdentities: &managedIdentitiesClient,
	}
}

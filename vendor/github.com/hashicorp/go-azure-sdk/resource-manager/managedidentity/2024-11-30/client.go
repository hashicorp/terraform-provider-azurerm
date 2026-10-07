package v2024_11_30

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/identities"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/systemassignedidentities"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	FederatedIdentityCredentials *federatedidentitycredentials.FederatedIdentityCredentialsClient
	Identities                   *identities.IdentitiesClient
	SystemAssignedIdentities     *systemassignedidentities.SystemAssignedIdentitiesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	federatedIdentityCredentialsClient, err := federatedidentitycredentials.NewFederatedIdentityCredentialsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FederatedIdentityCredentials client: %+v", err)
	}
	configureFunc(federatedIdentityCredentialsClient.Client)

	identitiesClient, err := identities.NewIdentitiesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Identities client: %+v", err)
	}
	configureFunc(identitiesClient.Client)

	systemAssignedIdentitiesClient, err := systemassignedidentities.NewSystemAssignedIdentitiesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SystemAssignedIdentities client: %+v", err)
	}
	configureFunc(systemAssignedIdentitiesClient.Client)

	return &Client{
		FederatedIdentityCredentials: federatedIdentityCredentialsClient,
		Identities:                   identitiesClient,
		SystemAssignedIdentities:     systemAssignedIdentitiesClient,
	}, nil
}

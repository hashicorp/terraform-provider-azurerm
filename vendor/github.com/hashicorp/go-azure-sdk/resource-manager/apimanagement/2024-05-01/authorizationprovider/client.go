package authorizationprovider

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationProviderClient struct {
	Client *resourcemanager.Client
}

func NewAuthorizationProviderClientWithBaseURI(sdkApi sdkEnv.Api) (*AuthorizationProviderClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "authorizationprovider", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AuthorizationProviderClient: %+v", err)
	}

	return &AuthorizationProviderClient{
		Client: client,
	}, nil
}

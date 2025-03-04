package authorizationrulesnamespaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationRulesNamespacesClient struct {
	Client *resourcemanager.Client
}

func NewAuthorizationRulesNamespacesClientWithBaseURI(sdkApi sdkEnv.Api) (*AuthorizationRulesNamespacesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "authorizationrulesnamespaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AuthorizationRulesNamespacesClient: %+v", err)
	}

	return &AuthorizationRulesNamespacesClient{
		Client: client,
	}, nil
}

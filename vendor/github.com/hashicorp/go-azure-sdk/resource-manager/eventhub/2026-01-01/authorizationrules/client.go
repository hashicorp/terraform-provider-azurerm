package authorizationrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationRulesClient struct {
	Client *resourcemanager.Client
}

func NewAuthorizationRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*AuthorizationRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "authorizationrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AuthorizationRulesClient: %+v", err)
	}

	return &AuthorizationRulesClient{
		Client: client,
	}, nil
}

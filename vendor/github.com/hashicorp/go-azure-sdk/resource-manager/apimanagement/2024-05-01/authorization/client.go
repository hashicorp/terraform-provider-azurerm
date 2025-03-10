package authorization

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationClient struct {
	Client *resourcemanager.Client
}

func NewAuthorizationClientWithBaseURI(sdkApi sdkEnv.Api) (*AuthorizationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "authorization", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AuthorizationClient: %+v", err)
	}

	return &AuthorizationClient{
		Client: client,
	}, nil
}

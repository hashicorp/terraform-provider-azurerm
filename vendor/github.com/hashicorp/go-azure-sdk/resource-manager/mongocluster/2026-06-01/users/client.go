package users

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsersClient struct {
	Client *resourcemanager.Client
}

func NewUsersClientWithBaseURI(sdkApi sdkEnv.Api) (*UsersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "users", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating UsersClient: %+v", err)
	}

	return &UsersClient{
		Client: client,
	}, nil
}

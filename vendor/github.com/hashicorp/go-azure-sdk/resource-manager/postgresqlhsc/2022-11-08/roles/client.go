package roles

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RolesClient struct {
	Client *resourcemanager.Client
}

func NewRolesClientWithBaseURI(sdkApi sdkEnv.Api) (*RolesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "roles", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RolesClient: %+v", err)
	}

	return &RolesClient{
		Client: client,
	}, nil
}

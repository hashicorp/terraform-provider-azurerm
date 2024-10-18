package service

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceClient struct {
	Client *resourcemanager.Client
}

func NewServiceClientWithBaseURI(sdkApi sdkEnv.Api) (*ServiceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "service", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServiceClient: %+v", err)
	}

	return &ServiceClient{
		Client: client,
	}, nil
}

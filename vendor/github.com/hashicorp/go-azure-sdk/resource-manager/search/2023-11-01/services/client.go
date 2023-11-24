package services

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesClient struct {
	Client *resourcemanager.Client
}

func NewServicesClientWithBaseURI(sdkApi sdkEnv.Api) (*ServicesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "services", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServicesClient: %+v", err)
	}

	return &ServicesClient{
		Client: client,
	}, nil
}

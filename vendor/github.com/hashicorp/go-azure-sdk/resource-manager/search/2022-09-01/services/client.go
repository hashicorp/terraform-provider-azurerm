package services

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesClient struct {
	Client *resourcemanager.Client
}

func NewServicesClientWithBaseURI(api environments.Api) (*ServicesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "services", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServicesClient: %+v", err)
	}

	return &ServicesClient{
		Client: client,
	}, nil
}

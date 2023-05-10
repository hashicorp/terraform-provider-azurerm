package providers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvidersClient struct {
	Client *resourcemanager.Client
}

func NewProvidersClientWithBaseURI(api environments.Api) (*ProvidersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "providers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProvidersClient: %+v", err)
	}

	return &ProvidersClient{
		Client: client,
	}, nil
}

package resourceproviders

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceProvidersClient struct {
	Client *resourcemanager.Client
}

func NewResourceProvidersClientWithBaseURI(sdkApi sdkEnv.Api) (*ResourceProvidersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "resourceproviders", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ResourceProvidersClient: %+v", err)
	}

	return &ResourceProvidersClient{
		Client: client,
	}, nil
}

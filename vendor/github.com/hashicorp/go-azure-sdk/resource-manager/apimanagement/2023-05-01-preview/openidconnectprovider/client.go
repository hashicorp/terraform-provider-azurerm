package openidconnectprovider

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenidConnectProviderClient struct {
	Client *resourcemanager.Client
}

func NewOpenidConnectProviderClientWithBaseURI(sdkApi sdkEnv.Api) (*OpenidConnectProviderClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "openidconnectprovider", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OpenidConnectProviderClient: %+v", err)
	}

	return &OpenidConnectProviderClient{
		Client: client,
	}, nil
}

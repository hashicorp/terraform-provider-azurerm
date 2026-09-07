package afdendpoints

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDEndpointsClient struct {
	Client *resourcemanager.Client
}

func NewAFDEndpointsClientWithBaseURI(sdkApi sdkEnv.Api) (*AFDEndpointsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "afdendpoints", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AFDEndpointsClient: %+v", err)
	}

	return &AFDEndpointsClient{
		Client: client,
	}, nil
}

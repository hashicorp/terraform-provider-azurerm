package afdorigingroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDOriginGroupsClient struct {
	Client *resourcemanager.Client
}

func NewAFDOriginGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*AFDOriginGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "afdorigingroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AFDOriginGroupsClient: %+v", err)
	}

	return &AFDOriginGroupsClient{
		Client: client,
	}, nil
}

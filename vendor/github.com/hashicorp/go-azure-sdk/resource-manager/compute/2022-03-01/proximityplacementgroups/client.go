package proximityplacementgroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProximityPlacementGroupsClient struct {
	Client *resourcemanager.Client
}

func NewProximityPlacementGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*ProximityPlacementGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "proximityplacementgroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProximityPlacementGroupsClient: %+v", err)
	}

	return &ProximityPlacementGroupsClient{
		Client: client,
	}, nil
}

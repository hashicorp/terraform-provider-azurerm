package trafficmanagergeographichierarchies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficManagerGeographicHierarchiesClient struct {
	Client *resourcemanager.Client
}

func NewTrafficManagerGeographicHierarchiesClientWithBaseURI(sdkApi sdkEnv.Api) (*TrafficManagerGeographicHierarchiesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "trafficmanagergeographichierarchies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TrafficManagerGeographicHierarchiesClient: %+v", err)
	}

	return &TrafficManagerGeographicHierarchiesClient{
		Client: client,
	}, nil
}

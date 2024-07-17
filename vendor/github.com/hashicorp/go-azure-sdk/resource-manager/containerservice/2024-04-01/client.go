package v2024_04_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetmembers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetupdatestrategies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/updateruns"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	FleetMembers          *fleetmembers.FleetMembersClient
	FleetUpdateStrategies *fleetupdatestrategies.FleetUpdateStrategiesClient
	Fleets                *fleets.FleetsClient
	UpdateRuns            *updateruns.UpdateRunsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	fleetMembersClient, err := fleetmembers.NewFleetMembersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FleetMembers client: %+v", err)
	}
	configureFunc(fleetMembersClient.Client)

	fleetUpdateStrategiesClient, err := fleetupdatestrategies.NewFleetUpdateStrategiesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FleetUpdateStrategies client: %+v", err)
	}
	configureFunc(fleetUpdateStrategiesClient.Client)

	fleetsClient, err := fleets.NewFleetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Fleets client: %+v", err)
	}
	configureFunc(fleetsClient.Client)

	updateRunsClient, err := updateruns.NewUpdateRunsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UpdateRuns client: %+v", err)
	}
	configureFunc(updateRunsClient.Client)

	return &Client{
		FleetMembers:          fleetMembersClient,
		FleetUpdateStrategies: fleetUpdateStrategiesClient,
		Fleets:                fleetsClient,
		UpdateRuns:            updateRunsClient,
	}, nil
}

package v2023_10_15

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/fleetmembers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/fleetupdatestrategies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/updateruns"
)

type Client struct {
	FleetMembers          *fleetmembers.FleetMembersClient
	FleetUpdateStrategies *fleetupdatestrategies.FleetUpdateStrategiesClient
	Fleets                *fleets.FleetsClient
	UpdateRuns            *updateruns.UpdateRunsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	fleetMembersClient := fleetmembers.NewFleetMembersClientWithBaseURI(endpoint)
	configureAuthFunc(&fleetMembersClient.Client)

	fleetUpdateStrategiesClient := fleetupdatestrategies.NewFleetUpdateStrategiesClientWithBaseURI(endpoint)
	configureAuthFunc(&fleetUpdateStrategiesClient.Client)

	fleetsClient := fleets.NewFleetsClientWithBaseURI(endpoint)
	configureAuthFunc(&fleetsClient.Client)

	updateRunsClient := updateruns.NewUpdateRunsClientWithBaseURI(endpoint)
	configureAuthFunc(&updateRunsClient.Client)

	return Client{
		FleetMembers:          &fleetMembersClient,
		FleetUpdateStrategies: &fleetUpdateStrategiesClient,
		Fleets:                &fleetsClient,
		UpdateRuns:            &updateRunsClient,
	}
}

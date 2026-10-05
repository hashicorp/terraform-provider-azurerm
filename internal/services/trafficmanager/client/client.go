// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagergeographichierarchies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	EndpointsClient              *trafficmanagers.TrafficmanagersClient
	GeographialHierarchiesClient *trafficmanagergeographichierarchies.TrafficManagerGeographicHierarchiesClient
	ProfilesClient               *profiles.ProfilesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	endpointsClient, err := trafficmanagers.NewTrafficmanagersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Endpoints Client: %+v", err)
	}
	o.Configure(endpointsClient.Client, o.Authorizers.ResourceManager)

	geographialHierarchiesClient, err := trafficmanagergeographichierarchies.NewTrafficManagerGeographicHierarchiesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Geographial Hierarchies Client: %+v", err)
	}
	o.Configure(geographialHierarchiesClient.Client, o.Authorizers.ResourceManager)

	profilesClient, err := profiles.NewProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Profiles Client: %+v", err)
	}
	o.Configure(profilesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		EndpointsClient:              endpointsClient,
		GeographialHierarchiesClient: geographialHierarchiesClient,
		ProfilesClient:               profilesClient,
	}, nil
}

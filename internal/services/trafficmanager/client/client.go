// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/geographichierarchies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	EndpointsClient              *endpoints.EndpointsClient
	GeographialHierarchiesClient *geographichierarchies.GeographicHierarchiesClient
	ProfilesClient               *profiles.ProfilesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	endpointsClient, err := endpoints.NewEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Endpoints Client: %+v", err)
	}
	o.Configure(endpointsClient.Client, o.Authorizers.ResourceManager)

	geographialHierarchiesClient, err := geographichierarchies.NewGeographicHierarchiesClientWithBaseURI(o.Environment.ResourceManager)
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

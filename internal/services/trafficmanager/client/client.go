package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/geographichierarchies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GeographialHierarchiesClient *geographichierarchies.GeographicHierarchiesClient
	EndpointsClient              *endpoints.EndpointsClient
	ProfilesClient               *profiles.ProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	endpointsClient := endpoints.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	geographialHierarchiesClient := geographichierarchies.NewGeographicHierarchiesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&geographialHierarchiesClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := profiles.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		EndpointsClient:              &endpointsClient,
		GeographialHierarchiesClient: &geographialHierarchiesClient,
		ProfilesClient:               &profilesClient,
	}
}

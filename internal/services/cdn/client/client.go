package client

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdcustomdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigins"
)

type Client struct {
	FrontDoorProfileEndpointsClient     *afdendpoints.AFDEndpointsClient
	FrontDoorProfileOriginGroupsClient  *afdorigingroups.AFDOriginGroupsClient
	FrontDoorProfileOriginsClient       *afdorigins.AFDOriginsClient
	FrontDoorProfileCustomDomainsClient *afdcustomdomains.AFDCustomDomainsClient
	CustomDomainsClient                 *cdn.CustomDomainsClient
	EndpointsClient                     *cdn.EndpointsClient
	ProfilesClient                      *cdn.ProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	frontDoorProfileEndpointsClient := afdendpoints.NewAFDEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorProfileEndpointsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorProfileOriginGroupsClient := afdorigingroups.NewAFDOriginGroupsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorProfileOriginGroupsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorProfileOriginsClient := afdorigins.NewAFDOriginsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorProfileOriginsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorProfileCustomDomainsClient := afdcustomdomains.NewAFDCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorProfileCustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := cdn.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdn.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdn.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FrontDoorProfileEndpointsClient:     &frontDoorProfileEndpointsClient,
		FrontDoorProfileOriginGroupsClient:  &frontDoorProfileOriginGroupsClient,
		FrontDoorProfileOriginsClient:       &frontDoorProfileOriginsClient,
		FrontDoorProfileCustomDomainsClient: &frontDoorProfileCustomDomainsClient,
		CustomDomainsClient:                 &customDomainsClient,
		EndpointsClient:                     &endpointsClient,
		ProfilesClient:                      &profilesClient,
	}
}

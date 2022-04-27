package client

import (
	cdnSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	cdnFrontDoorSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FrontDoorEndpointsClient        *cdnFrontDoorSdk.AFDEndpointsClient
	FrontDoorOriginGroupsClient     *cdnFrontDoorSdk.AFDOriginGroupsClient
	FrontDoorOriginsClient          *cdnFrontDoorSdk.AFDOriginsClient
	FrontDoorCustomDomainsClient    *cdnFrontDoorSdk.AFDCustomDomainsClient
	FrontdoorSecurityPoliciesClient *cdnFrontDoorSdk.SecurityPoliciesClient
	FrontdoorRoutesClient           *cdnFrontDoorSdk.RoutesClient
	FrontdoorRulesClient            *cdnFrontDoorSdk.RulesClient
	FrontdoorProfileClient          *cdnFrontDoorSdk.ProfilesClient
	FrontdoorSecretsClient          *cdnFrontDoorSdk.SecretsClient
	FrontdoorRuleSetsClient         *cdnFrontDoorSdk.RuleSetsClient
	FrontdoorLegacyPoliciesClient   *frontdoor.PoliciesClient
	CustomDomainsClient             *cdnSdk.CustomDomainsClient
	EndpointsClient                 *cdnSdk.EndpointsClient
	ProfilesClient                  *cdnSdk.ProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	frontDoorEndpointsClient := cdnFrontDoorSdk.NewAFDEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorEndpointsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginGroupsClient := cdnFrontDoorSdk.NewAFDOriginGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginGroupsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginsClient := cdnFrontDoorSdk.NewAFDOriginsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorCustomDomainsClient := cdnFrontDoorSdk.NewAFDCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorCustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	frontdoorPolicySecurityPoliciesClient := cdnFrontDoorSdk.NewSecurityPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorPolicySecurityPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorLegacyPoliciesClient := frontdoor.NewPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorLegacyPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRoutesClient := cdnFrontDoorSdk.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorRoutesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRulesClient := cdnFrontDoorSdk.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorRulesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorProfilesClient := cdnFrontDoorSdk.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorProfilesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorPolicySecretsClient := cdnFrontDoorSdk.NewSecretsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorPolicySecretsClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRuleSetsClient := cdnFrontDoorSdk.NewRuleSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorRuleSetsClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := cdnSdk.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdnSdk.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdnSdk.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FrontDoorEndpointsClient:        &frontDoorEndpointsClient,
		FrontDoorOriginGroupsClient:     &frontDoorOriginGroupsClient,
		FrontDoorOriginsClient:          &frontDoorOriginsClient,
		FrontDoorCustomDomainsClient:    &frontDoorCustomDomainsClient,
		FrontdoorSecurityPoliciesClient: &frontdoorPolicySecurityPoliciesClient,
		FrontdoorRoutesClient:           &frontdoorRoutesClient,
		FrontdoorRulesClient:            &frontdoorRulesClient,
		FrontdoorProfileClient:          &frontdoorProfilesClient,
		FrontdoorSecretsClient:          &frontdoorPolicySecretsClient,
		FrontdoorRuleSetsClient:         &frontdoorRuleSetsClient,
		FrontdoorLegacyPoliciesClient:   &frontdoorLegacyPoliciesClient,
		CustomDomainsClient:             &customDomainsClient,
		EndpointsClient:                 &endpointsClient,
		ProfilesClient:                  &profilesClient,
	}
}

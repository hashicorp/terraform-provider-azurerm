package client

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	legacyfrontdoor "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/legacysdk/2020-11-01"
	sdk "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
)

type Client struct {
	FrontDoorEndpointsClient        *sdk.AFDEndpointsClient
	FrontDoorOriginGroupsClient     *sdk.AFDOriginGroupsClient
	FrontDoorOriginsClient          *sdk.AFDOriginsClient
	FrontDoorCustomDomainsClient    *sdk.AFDCustomDomainsClient
	FrontdoorSecurityPoliciesClient *sdk.SecurityPoliciesClient
	FrontdoorRoutesClient           *sdk.RoutesClient
	FrontdoorRulesClient            *sdk.RulesClient
	FrontdoorProfileClient          *sdk.ProfilesClient
	FrontdoorSecretsClient          *sdk.SecretsClient
	FrontdoorRuleSetsClient         *sdk.RuleSetsClient
	FrontdoorLegacyPoliciesClient   *legacyfrontdoor.PoliciesClient
	CustomDomainsClient             *cdn.CustomDomainsClient
	EndpointsClient                 *cdn.EndpointsClient
	ProfilesClient                  *cdn.ProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	frontDoorEndpointsClient := sdk.NewAFDEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorEndpointsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginGroupsClient := sdk.NewAFDOriginGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginGroupsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginsClient := sdk.NewAFDOriginsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorCustomDomainsClient := sdk.NewAFDCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorCustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	frontdoorPolicySecurityPoliciesClient := sdk.NewSecurityPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorPolicySecurityPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorLegacyPoliciesClient := legacyfrontdoor.NewPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorLegacyPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRoutesClient := sdk.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorRoutesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRulesClient := sdk.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorRulesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorProfilesClient := sdk.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorProfilesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorPolicySecretsClient := sdk.NewSecretsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorPolicySecretsClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRuleSetsClient := sdk.NewRuleSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontdoorRuleSetsClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := cdn.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdn.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdn.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
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

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdcustomdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigins"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/webapplicationfirewallpolicies"
)

type Client struct {
	FrontDoorEndpointsClient             *afdendpoints.AFDEndpointsClient
	FrontDoorOriginGroupsClient          *afdorigingroups.AFDOriginGroupsClient
	FrontDoorOriginsClient               *afdorigins.AFDOriginsClient
	FrontDoorCustomDomainsClient         *afdcustomdomains.AFDCustomDomainsClient
	FrontdoorSecurityPoliciesClient      *securitypolicies.SecurityPoliciesClient
	FrontdoorRoutesClient                *routes.RoutesClient
	FrontdoorRulesClient                 *rules.RulesClient
	FrontdoorProfileClient               *profiles.ProfilesClient
	FrontdoorSecretsClient               *secrets.SecretsClient
	FrontdoorRuleSetsClient              *rulesets.RuleSetsClient
	WebApplicationFirewallPoliciesClient *webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient
	CustomDomainsClient                  *cdn.CustomDomainsClient
	EndpointsClient                      *cdn.EndpointsClient
	ProfilesClient                       *cdn.ProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	frontDoorEndpointsClient := afdendpoints.NewAFDEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorEndpointsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginGroupsClient := afdorigingroups.NewAFDOriginGroupsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorOriginGroupsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginsClient := afdorigins.NewAFDOriginsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorOriginsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorCustomDomainsClient := afdcustomdomains.NewAFDCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorCustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	frontdoorPolicySecurityPoliciesClient := securitypolicies.NewSecurityPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorPolicySecurityPoliciesClient.Client, o.ResourceManagerAuthorizer)

	webApplicationFirewallPoliciesClient := webapplicationfirewallpolicies.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&webApplicationFirewallPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRoutesClient := routes.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorRoutesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRulesClient := rules.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorRulesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorProfilesClient := profiles.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorProfilesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorPolicySecretsClient := secrets.NewSecretsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorPolicySecretsClient.Client, o.ResourceManagerAuthorizer)

	frontdoorRuleSetsClient := rulesets.NewRuleSetsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorRuleSetsClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := cdn.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdn.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdn.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FrontDoorEndpointsClient:             &frontDoorEndpointsClient,
		FrontDoorOriginGroupsClient:          &frontDoorOriginGroupsClient,
		FrontDoorOriginsClient:               &frontDoorOriginsClient,
		FrontDoorCustomDomainsClient:         &frontDoorCustomDomainsClient,
		FrontdoorSecurityPoliciesClient:      &frontdoorPolicySecurityPoliciesClient,
		FrontdoorRoutesClient:                &frontdoorRoutesClient,
		FrontdoorRulesClient:                 &frontdoorRulesClient,
		FrontdoorProfileClient:               &frontdoorProfilesClient,
		FrontdoorSecretsClient:               &frontdoorPolicySecretsClient,
		FrontdoorRuleSetsClient:              &frontdoorRuleSetsClient,
		WebApplicationFirewallPoliciesClient: &webApplicationFirewallPoliciesClient,
		CustomDomainsClient:                  &customDomainsClient,
		EndpointsClient:                      &endpointsClient,
		ProfilesClient:                       &profilesClient,
	}
}

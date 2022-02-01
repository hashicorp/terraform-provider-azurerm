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
	FrontDoorProfileEndpointsClient        *afdendpoints.AFDEndpointsClient
	FrontDoorProfileOriginGroupsClient     *afdorigingroups.AFDOriginGroupsClient
	FrontDoorProfileOriginsClient          *afdorigins.AFDOriginsClient
	FrontDoorProfileCustomDomainsClient    *afdcustomdomains.AFDCustomDomainsClient
	FrontdoorProfileSecurityPoliciesClient *securitypolicies.SecurityPoliciesClient
	FrontdoorProfileRoutesClient           *routes.RoutesClient
	FrontdoorProfileRulesClient            *rules.RulesClient
	FrontdoorProfileClient                 *profiles.ProfilesClient
	FrontdoorProfileSecretsClient          *secrets.SecretsClient
	FrontdoorProfileRuleSetsClient         *rulesets.RuleSetsClient
	WebApplicationFirewallPoliciesClient   *webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient
	CustomDomainsClient                    *cdn.CustomDomainsClient
	EndpointsClient                        *cdn.EndpointsClient
	ProfilesClient                         *cdn.ProfilesClient
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

	frontdoorProfilePolicySecurityPoliciesClient := securitypolicies.NewSecurityPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorProfilePolicySecurityPoliciesClient.Client, o.ResourceManagerAuthorizer)

	webApplicationFirewallPoliciesClient := webapplicationfirewallpolicies.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&webApplicationFirewallPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorProfileRoutesClient := routes.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorProfileRoutesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorProfileRulesClient := rules.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorProfileRulesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorProfilesClient := profiles.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorProfilesClient.Client, o.ResourceManagerAuthorizer)

	frontdoorProfilePolicySecretsClient := secrets.NewSecretsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorProfilePolicySecretsClient.Client, o.ResourceManagerAuthorizer)

	frontdoorProfileRuleSetsClient := rulesets.NewRuleSetsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontdoorProfileRuleSetsClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := cdn.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdn.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdn.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FrontDoorProfileEndpointsClient:        &frontDoorProfileEndpointsClient,
		FrontDoorProfileOriginGroupsClient:     &frontDoorProfileOriginGroupsClient,
		FrontDoorProfileOriginsClient:          &frontDoorProfileOriginsClient,
		FrontDoorProfileCustomDomainsClient:    &frontDoorProfileCustomDomainsClient,
		FrontdoorProfileSecurityPoliciesClient: &frontdoorProfilePolicySecurityPoliciesClient,
		FrontdoorProfileRoutesClient:           &frontdoorProfileRoutesClient,
		FrontdoorProfileRulesClient:            &frontdoorProfileRulesClient,
		FrontdoorProfileClient:                 &frontdoorProfilesClient,
		FrontdoorProfileSecretsClient:          &frontdoorProfilePolicySecretsClient,
		FrontdoorProfileRuleSetsClient:         &frontdoorProfileRuleSetsClient,
		WebApplicationFirewallPoliciesClient:   &webApplicationFirewallPoliciesClient,
		CustomDomainsClient:                    &customDomainsClient,
		EndpointsClient:                        &endpointsClient,
		ProfilesClient:                         &profilesClient,
	}
}

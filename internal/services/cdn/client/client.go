// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	cdnSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/secrets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/securitypolicies"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AFDCustomDomainsClient          *afddomains.AFDDomainsClient
	AFDEndpointsClient              *afdendpoints.AFDEndpointsClient
	FrontDoorFirewallPoliciesClient *waf.WebApplicationFirewallPoliciesClient
	FrontDoorOriginGroupsClient     *afdorigingroups.AFDOriginGroupsClient
	FrontDoorOriginsClient          *afdorigins.AFDOriginsClient
	FrontDoorProfilesClient         *profiles.ProfilesClient
	FrontDoorRoutesClient           *routes.RoutesClient
	FrontDoorRulesClient            *rules.RulesClient
	FrontDoorRuleSetsClient         *rulesets.RuleSetsClient
	FrontDoorSecretsClient          *secrets.SecretsClient
	FrontDoorSecurityPoliciesClient *securitypolicies.SecurityPoliciesClient

	// These clients are in-use by deprecated resources/data sources, and can no longer be created.
	// Because we are unable to test, we'll leave these on the track1 SDK.
	CustomDomainsClient *cdnSdk.CustomDomainsClient
	EndpointsClient     *cdnSdk.EndpointsClient
	ProfilesClient      *cdnSdk.ProfilesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	afdCustomDomainsClient, err := afddomains.NewAFDDomainsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AFD Domains Client: %+v", err)
	}
	o.Configure(afdCustomDomainsClient.Client, o.Authorizers.ResourceManager)

	afdEndpointsClient, err := afdendpoints.NewAFDEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Azure Front Door Endpoints CLient: %+v", err)
	}
	o.Configure(afdEndpointsClient.Client, o.Authorizers.ResourceManager)

	frontDoorFirewallPoliciesClient := waf.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorFirewallPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginGroupsClient, err := afdorigingroups.NewAFDOriginGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AFD Origin Groups Client: %+v", err)
	}
	o.Configure(frontDoorOriginGroupsClient.Client, o.Authorizers.ResourceManager)

	frontDoorOriginsClient, err := afdorigins.NewAFDOriginsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AFD Origins Client: %+v", err)
	}
	o.Configure(frontDoorOriginsClient.Client, o.Authorizers.ResourceManager)

	frontDoorProfilesClient, err := profiles.NewProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Profiles Client: %+v", err)
	}
	o.Configure(frontDoorProfilesClient.Client, o.Authorizers.ResourceManager)

	frontDoorRoutesClient, err := routes.NewRoutesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Routes Client: %+v", err)
	}
	o.Configure(frontDoorRoutesClient.Client, o.Authorizers.ResourceManager)

	frontDoorRulesClient, err := rules.NewRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Rules Client: %+v", err)
	}
	o.Configure(frontDoorRulesClient.Client, o.Authorizers.ResourceManager)

	frontDoorRuleSetsClient, err := rulesets.NewRuleSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Rule Sets Client: %+v", err)
	}
	o.Configure(frontDoorRuleSetsClient.Client, o.Authorizers.ResourceManager)

	frontDoorPolicySecretsClient, err := secrets.NewSecretsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Secrets Client: %+v", err)
	}
	o.Configure(frontDoorPolicySecretsClient.Client, o.Authorizers.ResourceManager)

	frontDoorSecurityPoliciesClient, err := securitypolicies.NewSecurityPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Security Policies Client: %+v", err)
	}
	o.Configure(frontDoorSecurityPoliciesClient.Client, o.Authorizers.ResourceManager)

	customDomainsClient := cdnSdk.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdnSdk.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdnSdk.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	client := Client{
		AFDCustomDomainsClient:          afdCustomDomainsClient,
		AFDEndpointsClient:              afdEndpointsClient,
		FrontDoorFirewallPoliciesClient: &frontDoorFirewallPoliciesClient,
		FrontDoorOriginGroupsClient:     frontDoorOriginGroupsClient,
		FrontDoorOriginsClient:          frontDoorOriginsClient,
		FrontDoorProfilesClient:         frontDoorProfilesClient,
		FrontDoorRoutesClient:           frontDoorRoutesClient,
		FrontDoorRulesClient:            frontDoorRulesClient,
		FrontDoorRuleSetsClient:         frontDoorRuleSetsClient,
		FrontDoorSecretsClient:          frontDoorPolicySecretsClient,
		FrontDoorSecurityPoliciesClient: frontDoorSecurityPoliciesClient,

		// Leave on track1 SDK
		CustomDomainsClient: &customDomainsClient,
		EndpointsClient:     &endpointsClient,
		ProfilesClient:      &profilesClient,
	}

	return &client, nil
}

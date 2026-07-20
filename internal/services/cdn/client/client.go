// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	cdnSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"          // nolint: staticcheck
	cdnFrontDoorSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/securitypolicies"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AFDCustomDomainsClient          *afddomains.AFDDomainsClient
	AFDEndpointsClient              *afdendpoints.AFDEndpointsClient
	FrontDoorSecurityPoliciesClient *securitypolicies.SecurityPoliciesClient
	FrontDoorRulesClient            *rules.RulesClient
	FrontDoorProfilesClient         *profiles.ProfilesClient
	FrontDoorRuleSetsClient         *rulesets.RuleSetsClient
	FrontDoorFirewallPoliciesClient *waf.WebApplicationFirewallPoliciesClient

	// TODO: migrate to go-azure-sdk
	FrontDoorOriginGroupsClient  *cdnFrontDoorSdk.AFDOriginGroupsClient
	FrontDoorOriginsClient       *cdnFrontDoorSdk.AFDOriginsClient
	FrontDoorCustomDomainsClient *cdnFrontDoorSdk.AFDCustomDomainsClient
	FrontDoorRoutesClient        *cdnFrontDoorSdk.RoutesClient
	FrontDoorSecretsClient       *cdnFrontDoorSdk.SecretsClient
	CustomDomainsClient          *cdnSdk.CustomDomainsClient
	EndpointsClient              *cdnSdk.EndpointsClient
	ProfilesClient               *cdnSdk.ProfilesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	frontDoorOriginGroupsClient := cdnFrontDoorSdk.NewAFDOriginGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginGroupsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginsClient := cdnFrontDoorSdk.NewAFDOriginsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorCustomDomainsClient := cdnFrontDoorSdk.NewAFDCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorCustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	afdCustomDomainsClient, err := afddomains.NewAFDDomainsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AFD Domains Client: %+v", err)
	}
	o.Configure(afdCustomDomainsClient.Client, o.Authorizers.ResourceManager)

	frontDoorSecurityPoliciesClient, err := securitypolicies.NewSecurityPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Security Policies Client: %+v", err)
	}
	o.Configure(frontDoorSecurityPoliciesClient.Client, o.Authorizers.ResourceManager)

	frontDoorFirewallPoliciesClient := waf.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorFirewallPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRoutesClient := cdnFrontDoorSdk.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorRoutesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRulesClient, err := rules.NewRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Rules Client: %+v", err)
	}
	o.Configure(frontDoorRulesClient.Client, o.Authorizers.ResourceManager)

	frontDoorProfilesClient, err := profiles.NewProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Profiles Client: %+v", err)
	}
	o.Configure(frontDoorProfilesClient.Client, o.Authorizers.ResourceManager)

	frontDoorPolicySecretsClient := cdnFrontDoorSdk.NewSecretsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorPolicySecretsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRuleSetsClient, err := rulesets.NewRuleSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Rule Sets Client: %+v", err)
	}
	o.Configure(frontDoorRuleSetsClient.Client, o.Authorizers.ResourceManager)

	customDomainsClient := cdnSdk.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdnSdk.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdnSdk.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	afdEndpointsClient, err := afdendpoints.NewAFDEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Azure Front Door Endpoints CLient: %+v", err)
	}
	o.Configure(afdEndpointsClient.Client, o.Authorizers.ResourceManager)

	client := Client{
		AFDCustomDomainsClient:          afdCustomDomainsClient,
		AFDEndpointsClient:              afdEndpointsClient,
		FrontDoorSecurityPoliciesClient: frontDoorSecurityPoliciesClient,
		FrontDoorRulesClient:            frontDoorRulesClient,
		FrontDoorProfilesClient:         frontDoorProfilesClient,
		FrontDoorRuleSetsClient:         frontDoorRuleSetsClient,
		FrontDoorFirewallPoliciesClient: &frontDoorFirewallPoliciesClient,

		// TODO: migrate to go-azure-sdk
		FrontDoorOriginGroupsClient:  &frontDoorOriginGroupsClient,
		FrontDoorOriginsClient:       &frontDoorOriginsClient,
		FrontDoorCustomDomainsClient: &frontDoorCustomDomainsClient,
		FrontDoorRoutesClient:        &frontDoorRoutesClient,
		FrontDoorSecretsClient:       &frontDoorPolicySecretsClient,
		CustomDomainsClient:          &customDomainsClient,
		EndpointsClient:              &endpointsClient,
		ProfilesClient:               &profilesClient,
	}

	return &client, nil
}

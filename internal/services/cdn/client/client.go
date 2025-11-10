// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	cdnSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"          // nolint: staticcheck
	cdnFrontDoorSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/securitypolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FrontDoorEndpointsClient        *cdnFrontDoorSdk.AFDEndpointsClient
	FrontDoorOriginGroupsClient     *cdnFrontDoorSdk.AFDOriginGroupsClient
	FrontDoorOriginsClient          *cdnFrontDoorSdk.AFDOriginsClient
	FrontDoorCustomDomainsClient    *cdnFrontDoorSdk.AFDCustomDomainsClient
	AFDCustomDomainsClient          *afdcustomdomains.AFDCustomDomainsClient
	FrontDoorSecurityPoliciesClient *securitypolicies.SecurityPoliciesClient
	FrontDoorRoutesClient           *cdnFrontDoorSdk.RoutesClient
	FrontDoorRulesClient            *rules.RulesClient
	FrontDoorProfilesClient         *profiles.ProfilesClient
	FrontDoorSecretsClient          *cdnFrontDoorSdk.SecretsClient
	FrontDoorRuleSetsClient         *rulesets.RuleSetsClient
	FrontDoorFirewallPoliciesClient *waf.WebApplicationFirewallPoliciesClient
	CustomDomainsClient             *cdnSdk.CustomDomainsClient
	EndpointsClient                 *cdnSdk.EndpointsClient
	ProfilesClient                  *cdnSdk.ProfilesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	frontDoorEndpointsClient := cdnFrontDoorSdk.NewAFDEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorEndpointsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginGroupsClient := cdnFrontDoorSdk.NewAFDOriginGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginGroupsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginsClient := cdnFrontDoorSdk.NewAFDOriginsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorCustomDomainsClient := cdnFrontDoorSdk.NewAFDCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorCustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	afdCustomDomainsClient, err := afdcustomdomains.NewAFDCustomDomainsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AFDCustomDomainsClient: %+v", err)
	}
	o.Configure(afdCustomDomainsClient.Client, o.Authorizers.ResourceManager)

	frontDoorSecurityPoliciesClient, err := securitypolicies.NewSecurityPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SecurityPoliciesClient: %+v", err)
	}
	o.Configure(frontDoorSecurityPoliciesClient.Client, o.Authorizers.ResourceManager)

	frontDoorFirewallPoliciesClient := waf.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorFirewallPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRoutesClient := cdnFrontDoorSdk.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorRoutesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRulesClient, err := rules.NewRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RulesClient: %+v", err)
	}
	o.Configure(frontDoorRulesClient.Client, o.Authorizers.ResourceManager)

	frontDoorProfilesClient, err := profiles.NewProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ProfilesClient: %+v", err)
	}
	o.Configure(frontDoorProfilesClient.Client, o.Authorizers.ResourceManager)

	frontDoorPolicySecretsClient := cdnFrontDoorSdk.NewSecretsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorPolicySecretsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRuleSetsClient, err := rulesets.NewRuleSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RuleSet: %+v", err)
	}
	o.Configure(frontDoorRuleSetsClient.Client, o.Authorizers.ResourceManager)

	customDomainsClient := cdnSdk.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdnSdk.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdnSdk.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	client := Client{
		FrontDoorEndpointsClient:        &frontDoorEndpointsClient,
		FrontDoorOriginGroupsClient:     &frontDoorOriginGroupsClient,
		FrontDoorOriginsClient:          &frontDoorOriginsClient,
		FrontDoorCustomDomainsClient:    &frontDoorCustomDomainsClient,
		AFDCustomDomainsClient:          afdCustomDomainsClient,
		FrontDoorSecurityPoliciesClient: frontDoorSecurityPoliciesClient,
		FrontDoorRoutesClient:           &frontDoorRoutesClient,
		FrontDoorRulesClient:            frontDoorRulesClient,
		FrontDoorProfilesClient:         frontDoorProfilesClient,
		FrontDoorSecretsClient:          &frontDoorPolicySecretsClient,
		FrontDoorRuleSetsClient:         frontDoorRuleSetsClient,
		FrontDoorFirewallPoliciesClient: &frontDoorFirewallPoliciesClient,
		CustomDomainsClient:             &customDomainsClient,
		EndpointsClient:                 &endpointsClient,
		ProfilesClient:                  &profilesClient,
	}

	return &client, nil
}

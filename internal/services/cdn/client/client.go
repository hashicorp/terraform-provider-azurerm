// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	cdnSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"          // nolint: staticcheck
	cdnFrontDoorSdk "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor"     // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FrontDoorEndpointsClient              *cdnFrontDoorSdk.AFDEndpointsClient
	FrontDoorOriginGroupsClient           *cdnFrontDoorSdk.AFDOriginGroupsClient
	FrontDoorOriginsClient                *cdnFrontDoorSdk.AFDOriginsClient
	FrontDoorCustomDomainsClient          *cdnFrontDoorSdk.AFDCustomDomainsClient
	FrontDoorSecurityPoliciesClient       *securitypolicies.SecurityPoliciesClient
	FrontDoorRoutesClient                 *cdnFrontDoorSdk.RoutesClient
	FrontDoorRulesClient                  *cdnFrontDoorSdk.RulesClient
	FrontDoorProfilesClient               *profiles.ProfilesClient
	FrontDoorSecretsClient                *cdnFrontDoorSdk.SecretsClient
	FrontDoorRuleSetsClient               *cdnFrontDoorSdk.RuleSetsClient
	FrontDoorLegacyFirewallPoliciesClient *frontdoor.PoliciesClient
	CustomDomainsClient                   *cdnSdk.CustomDomainsClient
	EndpointsClient                       *cdnSdk.EndpointsClient
	ProfilesClient                        *cdnSdk.ProfilesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	frontDoorProfilesClient, err := profiles.NewProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ProfilesClient: %+v", err)
	}
	o.Configure(frontDoorProfilesClient.Client, o.Authorizers.ResourceManager)

	frontDoorEndpointsClient := cdnFrontDoorSdk.NewAFDEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorEndpointsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginGroupsClient := cdnFrontDoorSdk.NewAFDOriginGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginGroupsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorOriginsClient := cdnFrontDoorSdk.NewAFDOriginsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorOriginsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorCustomDomainsClient := cdnFrontDoorSdk.NewAFDCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorCustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorSecurityPoliciesClient, err := securitypolicies.NewSecurityPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SecurityPoliciesClient: %+v", err)
	}
	o.Configure(frontDoorSecurityPoliciesClient.Client, o.Authorizers.ResourceManager)

	frontDoorLegacyFirewallPoliciesClient := frontdoor.NewPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorLegacyFirewallPoliciesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRoutesClient := cdnFrontDoorSdk.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorRoutesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRulesClient := cdnFrontDoorSdk.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorRulesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorPolicySecretsClient := cdnFrontDoorSdk.NewSecretsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorPolicySecretsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorRuleSetsClient := cdnFrontDoorSdk.NewRuleSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorRuleSetsClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := cdnSdk.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdnSdk.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdnSdk.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	client := Client{
		FrontDoorEndpointsClient:              &frontDoorEndpointsClient,
		FrontDoorOriginGroupsClient:           &frontDoorOriginGroupsClient,
		FrontDoorOriginsClient:                &frontDoorOriginsClient,
		FrontDoorCustomDomainsClient:          &frontDoorCustomDomainsClient,
		FrontDoorSecurityPoliciesClient:       frontDoorSecurityPoliciesClient,
		FrontDoorRoutesClient:                 &frontDoorRoutesClient,
		FrontDoorRulesClient:                  &frontDoorRulesClient,
		FrontDoorProfilesClient:               frontDoorProfilesClient,
		FrontDoorSecretsClient:                &frontDoorPolicySecretsClient,
		FrontDoorRuleSetsClient:               &frontDoorRuleSetsClient,
		FrontDoorLegacyFirewallPoliciesClient: &frontDoorLegacyFirewallPoliciesClient,
		CustomDomainsClient:                   &customDomainsClient,
		EndpointsClient:                       &endpointsClient,
		ProfilesClient:                        &profilesClient,
	}

	return &client, nil
}

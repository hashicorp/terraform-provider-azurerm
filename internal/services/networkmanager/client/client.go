// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/adminrulecollections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/adminrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/connectivityconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/networkgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/networkmanagerconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/networkmanagerroutingconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/networkmanagers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/networksecurityperimeterassociations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/networksecurityperimeters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/reachabilityanalysisintent"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/reachabilityanalysisintents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/routingrulecollections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/routingrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/scopeconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/securityadminconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/staticcidrs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/staticmembers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/verifierworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AdminRuleCollections                       *adminrulecollections.AdminRuleCollectionsClient
	AdminRules                                 *adminrules.AdminRulesClient
	ConnectivityConfigurations                 *connectivityconfigurations.ConnectivityConfigurationsClient
	IPamPools                                  *ipampools.IPamPoolsClient
	NetworkGroups                              *networkgroups.NetworkGroupsClient
	NetworkManagerConnections                  *networkmanagerconnections.NetworkManagerConnectionsClient
	NetworkManagerRoutingConfigurations        *networkmanagerroutingconfigurations.NetworkManagerRoutingConfigurationsClient
	NetworkManagers                            *networkmanagers.NetworkManagersClient
	NetworkSecurityPerimeterAccessRulesClient  *networksecurityperimeteraccessrules.NetworkSecurityPerimeterAccessRulesClient
	NetworkSecurityPerimeterAssociationsClient *networksecurityperimeterassociations.NetworkSecurityPerimeterAssociationsClient
	NetworkSecurityPerimeterProfilesClient     *networksecurityperimeterprofiles.NetworkSecurityPerimeterProfilesClient
	NetworkSecurityPerimetersClient            *networksecurityperimeters.NetworkSecurityPerimetersClient
	ReachabilityAnalysisIntent                 *reachabilityanalysisintent.ReachabilityAnalysisIntentClient
	ReachabilityAnalysisIntents                *reachabilityanalysisintents.ReachabilityAnalysisIntentsClient
	RoutingRuleCollections                     *routingrulecollections.RoutingRuleCollectionsClient
	RoutingRules                               *routingrules.RoutingRulesClient
	ScopeConnections                           *scopeconnections.ScopeConnectionsClient
	SecurityAdminConfigurations                *securityadminconfigurations.SecurityAdminConfigurationsClient
	StaticCidrs                                *staticcidrs.StaticCidrsClient
	StaticMembers                              *staticmembers.StaticMembersClient
	VerifierWorkspaces                         *verifierworkspaces.VerifierWorkspacesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	adminRuleCollectionsClient, err := adminrulecollections.NewAdminRuleCollectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Admin Rule Collections Client: %+v", err)
	}
	o.Configure(adminRuleCollectionsClient.Client, o.Authorizers.ResourceManager)

	adminRulesClient, err := adminrules.NewAdminRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Admin Rules Client: %+v", err)
	}
	o.Configure(adminRulesClient.Client, o.Authorizers.ResourceManager)

	connectivityConfigurationsClient, err := connectivityconfigurations.NewConnectivityConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Connectivity Configurations Client: %+v", err)
	}
	o.Configure(connectivityConfigurationsClient.Client, o.Authorizers.ResourceManager)

	ipamPoolsClient, err := ipampools.NewIPamPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building IPAM Pools Client: %+v", err)
	}
	o.Configure(ipamPoolsClient.Client, o.Authorizers.ResourceManager)

	networkGroupsClient, err := networkgroups.NewNetworkGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Groups Client: %+v", err)
	}
	o.Configure(networkGroupsClient.Client, o.Authorizers.ResourceManager)

	networkManagerConnectionsClient, err := networkmanagerconnections.NewNetworkManagerConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Manager Connections Client: %+v", err)
	}
	o.Configure(networkManagerConnectionsClient.Client, o.Authorizers.ResourceManager)

	networkManagerRoutingConfigurationsClient, err := networkmanagerroutingconfigurations.NewNetworkManagerRoutingConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Manager Routing Configurations Client: %+v", err)
	}
	o.Configure(networkManagerRoutingConfigurationsClient.Client, o.Authorizers.ResourceManager)

	networkManagersClient, err := networkmanagers.NewNetworkManagersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Managers Client: %+v", err)
	}
	o.Configure(networkManagersClient.Client, o.Authorizers.ResourceManager)

	networkSecurityPerimeterAccessRulesClient, err := networksecurityperimeteraccessrules.NewNetworkSecurityPerimeterAccessRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Security Perimeter Access Rules Client: %+v", err)
	}
	o.Configure(networkSecurityPerimeterAccessRulesClient.Client, o.Authorizers.ResourceManager)

	networkSecurityPerimeterAssociationsClient, err := networksecurityperimeterassociations.NewNetworkSecurityPerimeterAssociationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Security Perimeter Resource Association Client: %+v", err)
	}
	o.Configure(networkSecurityPerimeterAssociationsClient.Client, o.Authorizers.ResourceManager)

	networkSecurityPerimeterProfilesClient, err := networksecurityperimeterprofiles.NewNetworkSecurityPerimeterProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Security Perimeter Profiles Client: %+v", err)
	}
	o.Configure(networkSecurityPerimeterProfilesClient.Client, o.Authorizers.ResourceManager)

	networkSecurityPerimetersClient, err := networksecurityperimeters.NewNetworkSecurityPerimetersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Security Perimeters Client: %+v", err)
	}
	o.Configure(networkSecurityPerimetersClient.Client, o.Authorizers.ResourceManager)

	reachabilityAnalysisIntentClient, err := reachabilityanalysisintent.NewReachabilityAnalysisIntentClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Reachability Analysis Intent Client: %+v", err)
	}
	o.Configure(reachabilityAnalysisIntentClient.Client, o.Authorizers.ResourceManager)

	reachabilityAnalysisIntentsClient, err := reachabilityanalysisintents.NewReachabilityAnalysisIntentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Reachability Analysis Intents Client: %+v", err)
	}
	o.Configure(reachabilityAnalysisIntentsClient.Client, o.Authorizers.ResourceManager)

	routingRuleCollectionsClient, err := routingrulecollections.NewRoutingRuleCollectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Routing Rule Collections Client: %+v", err)
	}
	o.Configure(routingRuleCollectionsClient.Client, o.Authorizers.ResourceManager)

	routingRulesClient, err := routingrules.NewRoutingRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Routing Rules Client: %+v", err)
	}
	o.Configure(routingRulesClient.Client, o.Authorizers.ResourceManager)

	scopeConnectionsClient, err := scopeconnections.NewScopeConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Scope Connections Client: %+v", err)
	}
	o.Configure(scopeConnectionsClient.Client, o.Authorizers.ResourceManager)

	securityAdminConfigurationsClient, err := securityadminconfigurations.NewSecurityAdminConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Security Admin Configurations Client: %+v", err)
	}
	o.Configure(securityAdminConfigurationsClient.Client, o.Authorizers.ResourceManager)

	staticCidrsClient, err := staticcidrs.NewStaticCidrsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Static Cidrs Client: %+v", err)
	}
	o.Configure(staticCidrsClient.Client, o.Authorizers.ResourceManager)

	staticMembersClient, err := staticmembers.NewStaticMembersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Static Members Client: %+v", err)
	}
	o.Configure(staticMembersClient.Client, o.Authorizers.ResourceManager)

	verifierWorkspacesClient, err := verifierworkspaces.NewVerifierWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Verifier Workspaces Client: %+v", err)
	}
	o.Configure(verifierWorkspacesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AdminRuleCollections:                       adminRuleCollectionsClient,
		AdminRules:                                 adminRulesClient,
		ConnectivityConfigurations:                 connectivityConfigurationsClient,
		IPamPools:                                  ipamPoolsClient,
		NetworkGroups:                              networkGroupsClient,
		NetworkManagerConnections:                  networkManagerConnectionsClient,
		NetworkManagerRoutingConfigurations:        networkManagerRoutingConfigurationsClient,
		NetworkManagers:                            networkManagersClient,
		NetworkSecurityPerimeterAccessRulesClient:  networkSecurityPerimeterAccessRulesClient,
		NetworkSecurityPerimeterAssociationsClient: networkSecurityPerimeterAssociationsClient,
		NetworkSecurityPerimeterProfilesClient:     networkSecurityPerimeterProfilesClient,
		NetworkSecurityPerimetersClient:            networkSecurityPerimetersClient,
		ReachabilityAnalysisIntent:                 reachabilityAnalysisIntentClient,
		ReachabilityAnalysisIntents:                reachabilityAnalysisIntentsClient,
		RoutingRuleCollections:                     routingRuleCollectionsClient,
		RoutingRules:                               routingRulesClient,
		ScopeConnections:                           scopeConnectionsClient,
		SecurityAdminConfigurations:                securityAdminConfigurationsClient,
		StaticCidrs:                                staticCidrsClient,
		StaticMembers:                              staticMembersClient,
		VerifierWorkspaces:                         verifierWorkspacesClient,
	}, nil
}

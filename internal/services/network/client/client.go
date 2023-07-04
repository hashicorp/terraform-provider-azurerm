package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/adminrulecollections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/adminrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/connectivityconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/networkgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/networkmanagerconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/networkmanagers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/routefilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/scopeconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/securityadminconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/securityrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-09-01/staticmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type Client struct {
	ApplicationGatewaysClient                *network.ApplicationGatewaysClient
	ApplicationSecurityGroupsClient          *network.ApplicationSecurityGroupsClient
	BastionHostsClient                       *network.BastionHostsClient
	ConfigurationPolicyGroupClient           *network.ConfigurationPolicyGroupsClient
	ConnectionMonitorsClient                 *network.ConnectionMonitorsClient
	DDOSProtectionPlansClient                *network.DdosProtectionPlansClient
	ExpressRouteAuthsClient                  *network.ExpressRouteCircuitAuthorizationsClient
	ExpressRouteCircuitsClient               *network.ExpressRouteCircuitsClient
	ExpressRouteCircuitConnectionClient      *network.ExpressRouteCircuitConnectionsClient
	ExpressRouteConnectionsClient            *network.ExpressRouteConnectionsClient
	ExpressRouteGatewaysClient               *network.ExpressRouteGatewaysClient
	ExpressRoutePeeringsClient               *network.ExpressRouteCircuitPeeringsClient
	ExpressRoutePortsClient                  *network.ExpressRoutePortsClient
	ExpressRoutePortAuthorizationsClient     *network.ExpressRoutePortAuthorizationsClient
	FlowLogsClient                           *network.FlowLogsClient
	HubRouteTableClient                      *network.HubRouteTablesClient
	HubVirtualNetworkConnectionClient        *network.HubVirtualNetworkConnectionsClient
	InterfacesClient                         *network.InterfacesClient
	IPGroupsClient                           *network.IPGroupsClient
	LocalNetworkGatewaysClient               *network.LocalNetworkGatewaysClient
	ManagersClient                           *networkmanagers.NetworkManagersClient
	ManagerAdminRulesClient                  *adminrules.AdminRulesClient
	ManagerAdminRuleCollectionsClient        *adminrulecollections.AdminRuleCollectionsClient
	ManagerConnectivityConfigurationsClient  *connectivityconfigurations.ConnectivityConfigurationsClient
	ManagerConnectionsClient                 *networkmanagerconnections.NetworkManagerConnectionsClient
	ManagerNetworkGroupsClient               *networkgroups.NetworkGroupsClient
	ManagerScopeConnectionsClient            *scopeconnections.ScopeConnectionsClient
	ManagerSecurityAdminConfigurationsClient *securityadminconfigurations.SecurityAdminConfigurationsClient
	ManagerStaticMembersClient               *staticmembers.StaticMembersClient
	NatRuleClient                            *network.NatRulesClient
	PointToSiteVpnGatewaysClient             *network.P2sVpnGatewaysClient
	ProfileClient                            *network.ProfilesClient
	PacketCapturesClient                     *network.PacketCapturesClient
	PrivateEndpointClient                    *network.PrivateEndpointsClient
	PublicIPsClient                          *network.PublicIPAddressesClient
	PublicIPPrefixesClient                   *network.PublicIPPrefixesClient
	RouteMapsClient                          *network.RouteMapsClient
	RoutesClient                             *routes.RoutesClient
	RouteFiltersClient                       *routefilters.RouteFiltersClient
	RouteTablesClient                        *routetables.RouteTablesClient
	SecurityGroupClient                      *network.SecurityGroupsClient
	SecurityPartnerProviderClient            *network.SecurityPartnerProvidersClient
	SecurityRuleClient                       *securityrules.SecurityRulesClient
	ServiceEndpointPoliciesClient            *network.ServiceEndpointPoliciesClient
	ServiceEndpointPolicyDefinitionsClient   *network.ServiceEndpointPolicyDefinitionsClient
	ServiceTagsClient                        *network.ServiceTagsClient
	SubnetsClient                            *network.SubnetsClient
	NatGatewayClient                         *network.NatGatewaysClient
	VirtualHubBgpConnectionClient            *network.VirtualHubBgpConnectionClient
	VirtualHubIPClient                       *network.VirtualHubIPConfigurationClient
	VnetGatewayConnectionsClient             *network.VirtualNetworkGatewayConnectionsClient
	VnetGatewayNatRuleClient                 *network.VirtualNetworkGatewayNatRulesClient
	VnetGatewayClient                        *network.VirtualNetworkGatewaysClient
	VnetClient                               *network.VirtualNetworksClient
	VnetPeeringsClient                       *network.VirtualNetworkPeeringsClient
	VirtualWanClient                         *network.VirtualWansClient
	VirtualHubClient                         *network.VirtualHubsClient
	VpnConnectionsClient                     *network.VpnConnectionsClient
	VpnGatewaysClient                        *network.VpnGatewaysClient
	VpnServerConfigurationsClient            *network.VpnServerConfigurationsClient
	VpnSitesClient                           *network.VpnSitesClient
	WatcherClient                            *network.WatchersClient
	WebApplicationFirewallPoliciesClient     *network.WebApplicationFirewallPoliciesClient
	PrivateDnsZoneGroupClient                *network.PrivateDNSZoneGroupsClient
	PrivateLinkServiceClient                 *network.PrivateLinkServicesClient
	ServiceAssociationLinkClient             *network.ServiceAssociationLinksClient
	ResourceNavigationLinkClient             *network.ResourceNavigationLinksClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	ApplicationGatewaysClient := network.NewApplicationGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationGatewaysClient.Client, o.ResourceManagerAuthorizer)

	ApplicationSecurityGroupsClient := network.NewApplicationSecurityGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationSecurityGroupsClient.Client, o.ResourceManagerAuthorizer)

	BastionHostsClient := network.NewBastionHostsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&BastionHostsClient.Client, o.ResourceManagerAuthorizer)

	configurationPolicyGroupClient := network.NewConfigurationPolicyGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationPolicyGroupClient.Client, o.ResourceManagerAuthorizer)

	ConnectionMonitorsClient := network.NewConnectionMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConnectionMonitorsClient.Client, o.ResourceManagerAuthorizer)

	DDOSProtectionPlansClient := network.NewDdosProtectionPlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DDOSProtectionPlansClient.Client, o.ResourceManagerAuthorizer)

	ExpressRouteAuthsClient := network.NewExpressRouteCircuitAuthorizationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRouteAuthsClient.Client, o.ResourceManagerAuthorizer)

	ExpressRouteCircuitsClient := network.NewExpressRouteCircuitsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRouteCircuitsClient.Client, o.ResourceManagerAuthorizer)

	ExpressRouteCircuitConnectionClient := network.NewExpressRouteCircuitConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRouteCircuitConnectionClient.Client, o.ResourceManagerAuthorizer)

	ExpressRouteConnectionsClient := network.NewExpressRouteConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRouteConnectionsClient.Client, o.ResourceManagerAuthorizer)

	ExpressRouteGatewaysClient := network.NewExpressRouteGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRouteGatewaysClient.Client, o.ResourceManagerAuthorizer)

	ExpressRoutePeeringsClient := network.NewExpressRouteCircuitPeeringsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRoutePeeringsClient.Client, o.ResourceManagerAuthorizer)

	ExpressRoutePortsClient := network.NewExpressRoutePortsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRoutePortsClient.Client, o.ResourceManagerAuthorizer)

	ExpressRoutePortAuthorizationsClient := network.NewExpressRoutePortAuthorizationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRoutePortAuthorizationsClient.Client, o.ResourceManagerAuthorizer)

	FlowLogsClient := network.NewFlowLogsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FlowLogsClient.Client, o.ResourceManagerAuthorizer)

	HubRouteTableClient := network.NewHubRouteTablesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HubRouteTableClient.Client, o.ResourceManagerAuthorizer)

	HubVirtualNetworkConnectionClient := network.NewHubVirtualNetworkConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HubVirtualNetworkConnectionClient.Client, o.ResourceManagerAuthorizer)

	InterfacesClient := network.NewInterfacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&InterfacesClient.Client, o.ResourceManagerAuthorizer)

	IpGroupsClient := network.NewIPGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&IpGroupsClient.Client, o.ResourceManagerAuthorizer)

	LocalNetworkGatewaysClient := network.NewLocalNetworkGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LocalNetworkGatewaysClient.Client, o.ResourceManagerAuthorizer)

	ManagersClient, err := networkmanagers.NewNetworkManagersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network managers client: %+v", err)
	}
	o.Configure(ManagersClient.Client, o.Authorizers.ResourceManager)

	ManagerAdminRulesClient, err := adminrules.NewAdminRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network manager admin rules client: %+v", err)
	}
	o.Configure(ManagerAdminRulesClient.Client, o.Authorizers.ResourceManager)

	ManagerAdminRuleCollectionsClient, err := adminrulecollections.NewAdminRuleCollectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network manager admin rule collections client: %+v", err)
	}
	o.Configure(ManagerAdminRuleCollectionsClient.Client, o.Authorizers.ResourceManager)

	ManagerConnectivityConfigurationsClient, err := connectivityconfigurations.NewConnectivityConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network manager connectivity configurations client: %+v", err)
	}
	o.Configure(ManagerConnectivityConfigurationsClient.Client, o.Authorizers.ResourceManager)

	ManagerConnectionsClient, err := networkmanagerconnections.NewNetworkManagerConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network maanger connections client: %+v", err)
	}
	o.Configure(ManagerConnectionsClient.Client, o.Authorizers.ResourceManager)

	ManagerNetworkGroupsClient, err := networkgroups.NewNetworkGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network manager network groups client: %+v", err)
	}
	o.Configure(ManagerNetworkGroupsClient.Client, o.Authorizers.ResourceManager)

	ManagerScopeConnectionsClient, err := scopeconnections.NewScopeConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network manager scope connections client: %+v", err)
	}
	o.Configure(ManagerScopeConnectionsClient.Client, o.Authorizers.ResourceManager)

	ManagerSecurityAdminConfigurationsClient, err := securityadminconfigurations.NewSecurityAdminConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network manager security admin configurations client: %+v", err)
	}
	o.Configure(ManagerSecurityAdminConfigurationsClient.Client, o.Authorizers.ResourceManager)

	ManagerStaticMembersClient, err := staticmembers.NewStaticMembersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network manager static members client: %+v", err)
	}
	o.Configure(ManagerStaticMembersClient.Client, o.Authorizers.ResourceManager)

	NatRuleClient := network.NewNatRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&NatRuleClient.Client, o.ResourceManagerAuthorizer)

	pointToSiteVpnGatewaysClient := network.NewP2sVpnGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&pointToSiteVpnGatewaysClient.Client, o.ResourceManagerAuthorizer)

	vpnServerConfigurationsClient := network.NewVpnServerConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vpnServerConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	ProfileClient := network.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProfileClient.Client, o.ResourceManagerAuthorizer)

	VnetClient := network.NewVirtualNetworksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VnetClient.Client, o.ResourceManagerAuthorizer)

	PacketCapturesClient := network.NewPacketCapturesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PacketCapturesClient.Client, o.ResourceManagerAuthorizer)

	PrivateEndpointClient := network.NewPrivateEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PrivateEndpointClient.Client, o.ResourceManagerAuthorizer)

	VnetPeeringsClient := network.NewVirtualNetworkPeeringsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VnetPeeringsClient.Client, o.ResourceManagerAuthorizer)

	PublicIPsClient := network.NewPublicIPAddressesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PublicIPsClient.Client, o.ResourceManagerAuthorizer)

	PublicIPPrefixesClient := network.NewPublicIPPrefixesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PublicIPPrefixesClient.Client, o.ResourceManagerAuthorizer)

	PrivateDnsZoneGroupClient := network.NewPrivateDNSZoneGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PrivateDnsZoneGroupClient.Client, o.ResourceManagerAuthorizer)

	PrivateLinkServiceClient := network.NewPrivateLinkServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PrivateLinkServiceClient.Client, o.ResourceManagerAuthorizer)

	RouteMapsClient := network.NewRouteMapsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RouteMapsClient.Client, o.ResourceManagerAuthorizer)

	RoutesClient, err := routes.NewRoutesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network routes client: %+v", err)
	}
	o.Configure(RoutesClient.Client, o.Authorizers.ResourceManager)

	RouteFiltersClient, err := routefilters.NewRouteFiltersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network route filters client: %+v", err)
	}
	o.Configure(RouteFiltersClient.Client, o.Authorizers.ResourceManager)

	RouteTablesClient, err := routetables.NewRouteTablesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network route tables client: %+v", err)
	}
	o.Configure(RouteTablesClient.Client, o.Authorizers.ResourceManager)

	SecurityGroupClient := network.NewSecurityGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SecurityGroupClient.Client, o.ResourceManagerAuthorizer)

	SecurityPartnerProviderClient := network.NewSecurityPartnerProvidersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SecurityPartnerProviderClient.Client, o.ResourceManagerAuthorizer)

	SecurityRuleClient, err := securityrules.NewSecurityRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building network security rule client: %+v", err)
	}
	o.Configure(SecurityRuleClient.Client, o.Authorizers.ResourceManager)

	ServiceEndpointPoliciesClient := network.NewServiceEndpointPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServiceEndpointPoliciesClient.Client, o.ResourceManagerAuthorizer)

	ServiceEndpointPolicyDefinitionsClient := network.NewServiceEndpointPolicyDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServiceEndpointPolicyDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	ServiceTagsClient := network.NewServiceTagsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServiceTagsClient.Client, o.ResourceManagerAuthorizer)

	SubnetsClient := network.NewSubnetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SubnetsClient.Client, o.ResourceManagerAuthorizer)

	VirtualHubBgpConnectionClient := network.NewVirtualHubBgpConnectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualHubBgpConnectionClient.Client, o.ResourceManagerAuthorizer)

	VirtualHubIPClient := network.NewVirtualHubIPConfigurationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualHubIPClient.Client, o.ResourceManagerAuthorizer)

	VnetGatewayClient := network.NewVirtualNetworkGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VnetGatewayClient.Client, o.ResourceManagerAuthorizer)

	NatGatewayClient := network.NewNatGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&NatGatewayClient.Client, o.ResourceManagerAuthorizer)

	VnetGatewayConnectionsClient := network.NewVirtualNetworkGatewayConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VnetGatewayConnectionsClient.Client, o.ResourceManagerAuthorizer)

	VnetGatewayNatRuleClient := network.NewVirtualNetworkGatewayNatRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VnetGatewayNatRuleClient.Client, o.ResourceManagerAuthorizer)

	VirtualWanClient := network.NewVirtualWansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualWanClient.Client, o.ResourceManagerAuthorizer)

	VirtualHubClient := network.NewVirtualHubsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualHubClient.Client, o.ResourceManagerAuthorizer)

	vpnGatewaysClient := network.NewVpnGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vpnGatewaysClient.Client, o.ResourceManagerAuthorizer)

	vpnConnectionsClient := network.NewVpnConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vpnConnectionsClient.Client, o.ResourceManagerAuthorizer)

	vpnSitesClient := network.NewVpnSitesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vpnSitesClient.Client, o.ResourceManagerAuthorizer)

	WatcherClient := network.NewWatchersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WatcherClient.Client, o.ResourceManagerAuthorizer)

	WebApplicationFirewallPoliciesClient := network.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WebApplicationFirewallPoliciesClient.Client, o.ResourceManagerAuthorizer)

	ServiceAssociationLinkClient := network.NewServiceAssociationLinksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServiceAssociationLinkClient.Client, o.ResourceManagerAuthorizer)

	ResourceNavigationLinkClient := network.NewResourceNavigationLinksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ResourceNavigationLinkClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApplicationGatewaysClient:                &ApplicationGatewaysClient,
		ApplicationSecurityGroupsClient:          &ApplicationSecurityGroupsClient,
		BastionHostsClient:                       &BastionHostsClient,
		ConfigurationPolicyGroupClient:           &configurationPolicyGroupClient,
		ConnectionMonitorsClient:                 &ConnectionMonitorsClient,
		DDOSProtectionPlansClient:                &DDOSProtectionPlansClient,
		ExpressRouteAuthsClient:                  &ExpressRouteAuthsClient,
		ExpressRouteCircuitsClient:               &ExpressRouteCircuitsClient,
		ExpressRouteCircuitConnectionClient:      &ExpressRouteCircuitConnectionClient,
		ExpressRouteConnectionsClient:            &ExpressRouteConnectionsClient,
		ExpressRouteGatewaysClient:               &ExpressRouteGatewaysClient,
		ExpressRoutePeeringsClient:               &ExpressRoutePeeringsClient,
		ExpressRoutePortsClient:                  &ExpressRoutePortsClient,
		ExpressRoutePortAuthorizationsClient:     &ExpressRoutePortAuthorizationsClient,
		FlowLogsClient:                           &FlowLogsClient,
		HubRouteTableClient:                      &HubRouteTableClient,
		HubVirtualNetworkConnectionClient:        &HubVirtualNetworkConnectionClient,
		InterfacesClient:                         &InterfacesClient,
		IPGroupsClient:                           &IpGroupsClient,
		LocalNetworkGatewaysClient:               &LocalNetworkGatewaysClient,
		ManagersClient:                           ManagersClient,
		ManagerAdminRulesClient:                  ManagerAdminRulesClient,
		ManagerAdminRuleCollectionsClient:        ManagerAdminRuleCollectionsClient,
		ManagerConnectivityConfigurationsClient:  ManagerConnectivityConfigurationsClient,
		ManagerConnectionsClient:                 ManagerConnectionsClient,
		ManagerNetworkGroupsClient:               ManagerNetworkGroupsClient,
		ManagerScopeConnectionsClient:            ManagerScopeConnectionsClient,
		ManagerSecurityAdminConfigurationsClient: ManagerSecurityAdminConfigurationsClient,
		ManagerStaticMembersClient:               ManagerStaticMembersClient,
		NatRuleClient:                            &NatRuleClient,
		PointToSiteVpnGatewaysClient:             &pointToSiteVpnGatewaysClient,
		ProfileClient:                            &ProfileClient,
		PacketCapturesClient:                     &PacketCapturesClient,
		PrivateEndpointClient:                    &PrivateEndpointClient,
		PublicIPsClient:                          &PublicIPsClient,
		PublicIPPrefixesClient:                   &PublicIPPrefixesClient,
		RouteMapsClient:                          &RouteMapsClient,
		RoutesClient:                             RoutesClient,
		RouteFiltersClient:                       RouteFiltersClient,
		RouteTablesClient:                        RouteTablesClient,
		SecurityGroupClient:                      &SecurityGroupClient,
		SecurityPartnerProviderClient:            &SecurityPartnerProviderClient,
		SecurityRuleClient:                       SecurityRuleClient,
		ServiceEndpointPoliciesClient:            &ServiceEndpointPoliciesClient,
		ServiceEndpointPolicyDefinitionsClient:   &ServiceEndpointPolicyDefinitionsClient,
		ServiceTagsClient:                        &ServiceTagsClient,
		SubnetsClient:                            &SubnetsClient,
		NatGatewayClient:                         &NatGatewayClient,
		VirtualHubBgpConnectionClient:            &VirtualHubBgpConnectionClient,
		VirtualHubIPClient:                       &VirtualHubIPClient,
		VnetGatewayConnectionsClient:             &VnetGatewayConnectionsClient,
		VnetGatewayNatRuleClient:                 &VnetGatewayNatRuleClient,
		VnetGatewayClient:                        &VnetGatewayClient,
		VnetClient:                               &VnetClient,
		VnetPeeringsClient:                       &VnetPeeringsClient,
		VirtualWanClient:                         &VirtualWanClient,
		VirtualHubClient:                         &VirtualHubClient,
		VpnConnectionsClient:                     &vpnConnectionsClient,
		VpnGatewaysClient:                        &vpnGatewaysClient,
		VpnServerConfigurationsClient:            &vpnServerConfigurationsClient,
		VpnSitesClient:                           &vpnSitesClient,
		WatcherClient:                            &WatcherClient,
		WebApplicationFirewallPoliciesClient:     &WebApplicationFirewallPoliciesClient,
		PrivateDnsZoneGroupClient:                &PrivateDnsZoneGroupClient,
		PrivateLinkServiceClient:                 &PrivateLinkServiceClient,
		ServiceAssociationLinkClient:             &ServiceAssociationLinkClient,
		ResourceNavigationLinkClient:             &ResourceNavigationLinkClient,
	}, nil
}

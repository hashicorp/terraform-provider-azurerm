package client

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationGatewaysClient              *network.ApplicationGatewaysClient
	ApplicationSecurityGroupsClient        *network.ApplicationSecurityGroupsClient
	BastionHostsClient                     *network.BastionHostsClient
	ConnectionMonitorsClient               *network.ConnectionMonitorsClient
	DDOSProtectionPlansClient              *network.DdosProtectionPlansClient
	ExpressRouteAuthsClient                *network.ExpressRouteCircuitAuthorizationsClient
	ExpressRouteCircuitsClient             *network.ExpressRouteCircuitsClient
	ExpressRouteGatewaysClient             *network.ExpressRouteGatewaysClient
	ExpressRoutePeeringsClient             *network.ExpressRouteCircuitPeeringsClient
	ExpressRoutePortsClient                *network.ExpressRoutePortsClient
	FlowLogsClient                         *network.FlowLogsClient
	HubRouteTableClient                    *network.HubRouteTablesClient
	HubVirtualNetworkConnectionClient      *network.HubVirtualNetworkConnectionsClient
	InterfacesClient                       *network.InterfacesClient
	IPGroupsClient                         *network.IPGroupsClient
	LocalNetworkGatewaysClient             *network.LocalNetworkGatewaysClient
	PointToSiteVpnGatewaysClient           *network.P2sVpnGatewaysClient
	ProfileClient                          *network.ProfilesClient
	PacketCapturesClient                   *network.PacketCapturesClient
	PrivateEndpointClient                  *network.PrivateEndpointsClient
	PublicIPsClient                        *network.PublicIPAddressesClient
	PublicIPPrefixesClient                 *network.PublicIPPrefixesClient
	RoutesClient                           *network.RoutesClient
	RouteFiltersClient                     *network.RouteFiltersClient
	RouteTablesClient                      *network.RouteTablesClient
	SecurityGroupClient                    *network.SecurityGroupsClient
	SecurityPartnerProviderClient          *network.SecurityPartnerProvidersClient
	SecurityRuleClient                     *network.SecurityRulesClient
	ServiceEndpointPoliciesClient          *network.ServiceEndpointPoliciesClient
	ServiceEndpointPolicyDefinitionsClient *network.ServiceEndpointPolicyDefinitionsClient
	ServiceTagsClient                      *network.ServiceTagsClient
	SubnetsClient                          *network.SubnetsClient
	NatGatewayClient                       *network.NatGatewaysClient
	VirtualHubBgpConnectionClient          *network.VirtualHubBgpConnectionClient
	VirtualHubIPClient                     *network.VirtualHubIPConfigurationClient
	VnetGatewayConnectionsClient           *network.VirtualNetworkGatewayConnectionsClient
	VnetGatewayClient                      *network.VirtualNetworkGatewaysClient
	VnetClient                             *network.VirtualNetworksClient
	VnetPeeringsClient                     *network.VirtualNetworkPeeringsClient
	VirtualWanClient                       *network.VirtualWansClient
	VirtualHubClient                       *network.VirtualHubsClient
	VpnConnectionsClient                   *network.VpnConnectionsClient
	VpnGatewaysClient                      *network.VpnGatewaysClient
	VpnServerConfigurationsClient          *network.VpnServerConfigurationsClient
	VpnSitesClient                         *network.VpnSitesClient
	WatcherClient                          *network.WatchersClient
	WebApplicationFirewallPoliciesClient   *network.WebApplicationFirewallPoliciesClient
	PrivateDnsZoneGroupClient              *network.PrivateDNSZoneGroupsClient
	PrivateLinkServiceClient               *network.PrivateLinkServicesClient
	ServiceAssociationLinkClient           *network.ServiceAssociationLinksClient
	ResourceNavigationLinkClient           *network.ResourceNavigationLinksClient
}

func NewClient(o *common.ClientOptions) *Client {
	ApplicationGatewaysClient := network.NewApplicationGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationGatewaysClient.Client, o.ResourceManagerAuthorizer)

	ApplicationSecurityGroupsClient := network.NewApplicationSecurityGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationSecurityGroupsClient.Client, o.ResourceManagerAuthorizer)

	BastionHostsClient := network.NewBastionHostsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&BastionHostsClient.Client, o.ResourceManagerAuthorizer)

	ConnectionMonitorsClient := network.NewConnectionMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConnectionMonitorsClient.Client, o.ResourceManagerAuthorizer)

	DDOSProtectionPlansClient := network.NewDdosProtectionPlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DDOSProtectionPlansClient.Client, o.ResourceManagerAuthorizer)

	ExpressRouteAuthsClient := network.NewExpressRouteCircuitAuthorizationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRouteAuthsClient.Client, o.ResourceManagerAuthorizer)

	ExpressRouteCircuitsClient := network.NewExpressRouteCircuitsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRouteCircuitsClient.Client, o.ResourceManagerAuthorizer)

	ExpressRouteGatewaysClient := network.NewExpressRouteGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRouteGatewaysClient.Client, o.ResourceManagerAuthorizer)

	ExpressRoutePeeringsClient := network.NewExpressRouteCircuitPeeringsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRoutePeeringsClient.Client, o.ResourceManagerAuthorizer)

	ExpressRoutePortsClient := network.NewExpressRoutePortsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExpressRoutePortsClient.Client, o.ResourceManagerAuthorizer)

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

	RoutesClient := network.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RoutesClient.Client, o.ResourceManagerAuthorizer)

	RouteFiltersClient := network.NewRouteFiltersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RouteFiltersClient.Client, o.ResourceManagerAuthorizer)

	RouteTablesClient := network.NewRouteTablesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RouteTablesClient.Client, o.ResourceManagerAuthorizer)

	SecurityGroupClient := network.NewSecurityGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SecurityGroupClient.Client, o.ResourceManagerAuthorizer)

	SecurityPartnerProviderClient := network.NewSecurityPartnerProvidersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SecurityPartnerProviderClient.Client, o.ResourceManagerAuthorizer)

	SecurityRuleClient := network.NewSecurityRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SecurityRuleClient.Client, o.ResourceManagerAuthorizer)

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
		ApplicationGatewaysClient:              &ApplicationGatewaysClient,
		ApplicationSecurityGroupsClient:        &ApplicationSecurityGroupsClient,
		BastionHostsClient:                     &BastionHostsClient,
		ConnectionMonitorsClient:               &ConnectionMonitorsClient,
		DDOSProtectionPlansClient:              &DDOSProtectionPlansClient,
		ExpressRouteAuthsClient:                &ExpressRouteAuthsClient,
		ExpressRouteCircuitsClient:             &ExpressRouteCircuitsClient,
		ExpressRouteGatewaysClient:             &ExpressRouteGatewaysClient,
		ExpressRoutePeeringsClient:             &ExpressRoutePeeringsClient,
		ExpressRoutePortsClient:                &ExpressRoutePortsClient,
		FlowLogsClient:                         &FlowLogsClient,
		HubRouteTableClient:                    &HubRouteTableClient,
		HubVirtualNetworkConnectionClient:      &HubVirtualNetworkConnectionClient,
		InterfacesClient:                       &InterfacesClient,
		IPGroupsClient:                         &IpGroupsClient,
		LocalNetworkGatewaysClient:             &LocalNetworkGatewaysClient,
		PointToSiteVpnGatewaysClient:           &pointToSiteVpnGatewaysClient,
		ProfileClient:                          &ProfileClient,
		PacketCapturesClient:                   &PacketCapturesClient,
		PrivateEndpointClient:                  &PrivateEndpointClient,
		PublicIPsClient:                        &PublicIPsClient,
		PublicIPPrefixesClient:                 &PublicIPPrefixesClient,
		RoutesClient:                           &RoutesClient,
		RouteFiltersClient:                     &RouteFiltersClient,
		RouteTablesClient:                      &RouteTablesClient,
		SecurityGroupClient:                    &SecurityGroupClient,
		SecurityPartnerProviderClient:          &SecurityPartnerProviderClient,
		SecurityRuleClient:                     &SecurityRuleClient,
		ServiceEndpointPoliciesClient:          &ServiceEndpointPoliciesClient,
		ServiceEndpointPolicyDefinitionsClient: &ServiceEndpointPolicyDefinitionsClient,
		ServiceTagsClient:                      &ServiceTagsClient,
		SubnetsClient:                          &SubnetsClient,
		NatGatewayClient:                       &NatGatewayClient,
		VirtualHubBgpConnectionClient:          &VirtualHubBgpConnectionClient,
		VirtualHubIPClient:                     &VirtualHubIPClient,
		VnetGatewayConnectionsClient:           &VnetGatewayConnectionsClient,
		VnetGatewayClient:                      &VnetGatewayClient,
		VnetClient:                             &VnetClient,
		VnetPeeringsClient:                     &VnetPeeringsClient,
		VirtualWanClient:                       &VirtualWanClient,
		VirtualHubClient:                       &VirtualHubClient,
		VpnConnectionsClient:                   &vpnConnectionsClient,
		VpnGatewaysClient:                      &vpnGatewaysClient,
		VpnServerConfigurationsClient:          &vpnServerConfigurationsClient,
		VpnSitesClient:                         &vpnSitesClient,
		WatcherClient:                          &WatcherClient,
		WebApplicationFirewallPoliciesClient:   &WebApplicationFirewallPoliciesClient,
		PrivateDnsZoneGroupClient:              &PrivateDnsZoneGroupClient,
		PrivateLinkServiceClient:               &PrivateLinkServiceClient,
		ServiceAssociationLinkClient:           &ServiceAssociationLinkClient,
		ResourceNavigationLinkClient:           &ResourceNavigationLinkClient,
	}
}

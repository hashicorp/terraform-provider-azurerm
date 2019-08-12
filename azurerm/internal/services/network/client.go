package network

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationGatewaysClient       network.ApplicationGatewaysClient
	ApplicationSecurityGroupsClient network.ApplicationSecurityGroupsClient
	AzureFirewallsClient            network.AzureFirewallsClient
	ConnectionMonitorsClient        network.ConnectionMonitorsClient
	DDOSProtectionPlansClient       network.DdosProtectionPlansClient
	ExpressRouteAuthsClient         network.ExpressRouteCircuitAuthorizationsClient
	ExpressRouteCircuitsClient      network.ExpressRouteCircuitsClient
	ExpressRoutePeeringsClient      network.ExpressRouteCircuitPeeringsClient
	InterfacesClient                network.InterfacesClient
	LoadBalancersClient             network.LoadBalancersClient
	LocalNetworkGatewaysClient      network.LocalNetworkGatewaysClient
	ProfileClient                   network.ProfilesClient
	PacketCapturesClient            network.PacketCapturesClient
	PublicIPsClient                 network.PublicIPAddressesClient
	PublicIPPrefixesClient          network.PublicIPPrefixesClient
	RoutesClient                    network.RoutesClient
	RouteTablesClient               network.RouteTablesClient
	SecurityGroupClient             network.SecurityGroupsClient
	SecurityRuleClient              network.SecurityRulesClient
	SubnetsClient                   network.SubnetsClient
	VnetGatewayConnectionsClient    network.VirtualNetworkGatewayConnectionsClient
	VnetGatewayClient               network.VirtualNetworkGatewaysClient
	VnetClient                      network.VirtualNetworksClient
	VnetPeeringsClient              network.VirtualNetworkPeeringsClient
	WatcherClient                   network.WatchersClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ApplicationGatewaysClient = network.NewApplicationGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApplicationGatewaysClient.Client, o.ResourceManagerAuthorizer)

	c.ApplicationSecurityGroupsClient = network.NewApplicationSecurityGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApplicationSecurityGroupsClient.Client, o.ResourceManagerAuthorizer)

	c.AzureFirewallsClient = network.NewAzureFirewallsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AzureFirewallsClient.Client, o.ResourceManagerAuthorizer)

	c.ConnectionMonitorsClient = network.NewConnectionMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ConnectionMonitorsClient.Client, o.ResourceManagerAuthorizer)

	c.DDOSProtectionPlansClient = network.NewDdosProtectionPlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DDOSProtectionPlansClient.Client, o.ResourceManagerAuthorizer)

	c.ExpressRouteAuthsClient = network.NewExpressRouteCircuitAuthorizationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ExpressRouteAuthsClient.Client, o.ResourceManagerAuthorizer)

	c.ExpressRouteCircuitsClient = network.NewExpressRouteCircuitsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ExpressRouteCircuitsClient.Client, o.ResourceManagerAuthorizer)

	c.ExpressRoutePeeringsClient = network.NewExpressRouteCircuitPeeringsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ExpressRoutePeeringsClient.Client, o.ResourceManagerAuthorizer)

	c.InterfacesClient = network.NewInterfacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.InterfacesClient.Client, o.ResourceManagerAuthorizer)

	c.LoadBalancersClient = network.NewLoadBalancersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LoadBalancersClient.Client, o.ResourceManagerAuthorizer)

	c.LocalNetworkGatewaysClient = network.NewLocalNetworkGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LocalNetworkGatewaysClient.Client, o.ResourceManagerAuthorizer)

	c.ProfileClient = network.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProfileClient.Client, o.ResourceManagerAuthorizer)

	c.VnetClient = network.NewVirtualNetworksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VnetClient.Client, o.ResourceManagerAuthorizer)

	c.PacketCapturesClient = network.NewPacketCapturesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PacketCapturesClient.Client, o.ResourceManagerAuthorizer)

	c.VnetPeeringsClient = network.NewVirtualNetworkPeeringsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VnetPeeringsClient.Client, o.ResourceManagerAuthorizer)

	c.PublicIPsClient = network.NewPublicIPAddressesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PublicIPsClient.Client, o.ResourceManagerAuthorizer)

	c.PublicIPPrefixesClient = network.NewPublicIPPrefixesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PublicIPPrefixesClient.Client, o.ResourceManagerAuthorizer)

	c.RoutesClient = network.NewRoutesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RoutesClient.Client, o.ResourceManagerAuthorizer)

	c.RouteTablesClient = network.NewRouteTablesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RouteTablesClient.Client, o.ResourceManagerAuthorizer)

	c.SecurityGroupClient = network.NewSecurityGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SecurityGroupClient.Client, o.ResourceManagerAuthorizer)

	c.SecurityRuleClient = network.NewSecurityRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SecurityRuleClient.Client, o.ResourceManagerAuthorizer)

	c.SubnetsClient = network.NewSubnetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SubnetsClient.Client, o.ResourceManagerAuthorizer)

	c.VnetGatewayClient = network.NewVirtualNetworkGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VnetGatewayClient.Client, o.ResourceManagerAuthorizer)

	c.VnetGatewayConnectionsClient = network.NewVirtualNetworkGatewayConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VnetGatewayConnectionsClient.Client, o.ResourceManagerAuthorizer)

	c.WatcherClient = network.NewWatchersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.WatcherClient.Client, o.ResourceManagerAuthorizer)

	return &c
}

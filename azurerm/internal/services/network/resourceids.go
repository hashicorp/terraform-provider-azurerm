package network

// Core bits and pieces
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/applicationGateway1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGatewayHTTPListener -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/applicationGateway1/httpListeners/httpListener1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGatewayURLPathMapPathRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/applicationGateway1/urlPathMaps/urlPathMap1/pathRules/pathRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IpGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/ipGroups/group1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkInterface -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkSecurityGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkSecurityGroups/securityGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PublicIpAddress -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/publicIPAddresses/publicIpAddress1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PublicIpPrefix -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/publicIPPrefixes/publicIpPrefix1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Route -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/routeTables/routeTable1/routes/route1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RouteTable -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/routeTables/routeTable1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Subnet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/subnets/subnet1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetwork -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGatewayWebApplicationFirewallPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies/applicationGatewayWebApplicationFirewallPolicy1

// Bastion
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BastionHost -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/bastionHosts/bastionHost1

// NAT Gateway
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NatGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/natGateways/gateway1
// NOTE: the Nat Gateway <-> Public IP Association can't be generated at this time

// Network Watcher
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ConnectionMonitor -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkWatchers/watcher1/connectionMonitors/connectionMonitor1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkWatcher -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkWatchers/watcher1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PacketCapture -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkWatchers/watcher1/packetCaptures/capture1

// Private Link
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateEndpoint -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateDnsZoneConfig -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1/privateDnsZoneConfigs/privateDnsZoneConfig1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateDnsZoneGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1
// ^ these two looks like it should be in Private DNS - alas no, it's actually nested and entirely managed within the Private Endpoint

// Routing
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RouteFilter -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/routeFilters/filter1

// Virtual Hubs
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BgpConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/bgpConnections/connection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HubRouteTable -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/routeTable1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HubVirtualNetworkConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubVirtualNetworkConnections/hubConnection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SecurityPartnerProvider -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/securityPartnerProviders/partnerProvider1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualHub -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualHubIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/ipConfigurations/ipConfiguration1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualWan -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualWans/virtualWan1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VpnGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PointToSiteVpnGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/p2sVpnGateways/pointToSite1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VpnConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1/vpnConnections/vpnConnection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VpnServerConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnServerConfigurations/serverConfiguration1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VpnSite -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnSites/vpnSite1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VpnSiteLink -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnSites/vpnSite1/vpnSiteLinks/vpnSiteLink1

// Subnet Service Endpoint Policy
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SubnetServiceEndpointStoragePolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/serviceEndpointPolicies/policy1

// Express Route Port
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRoutePort -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRoutePorts/port1

// Virtual Network Gateway
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkGatewayIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/ipConfigurations/cfg1

// Express Route Connection
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteCircuit -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/erCircuit1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteCircuitPeering -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/erCircuit1/peerings/peering1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteGateways/ergw1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteCircuitConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/connection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteGateways/ergw1/expressRouteConnections/erConnection1

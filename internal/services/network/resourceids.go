// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

// Core bits and pieces
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/applicationGateway1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGatewayHTTPListener -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/applicationGateway1/httpListeners/httpListener1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGatewayURLPathMapPathRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/applicationGateway1/urlPathMaps/urlPathMap1/pathRules/pathRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IpGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/ipGroups/group1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkInterface -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkSecurityGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkSecurityGroups/securityGroup1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PublicIpAddress -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/publicIPAddresses/publicIpAddress1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PublicIpPrefix -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/publicIPPrefixes/publicIpPrefix1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkDnsServers -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/dnsServers/default -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkPeering -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/vnet1/virtualNetworkPeerings/vnetPeering1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkGatewayConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/connections/connection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DdosProtectionPlan -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/ddosProtectionPlans/ddosProtectionPlan1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateLinkService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateLinkServices/privateLinkService1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LocalNetworkGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/localNetworkGateways/localNetworkGateway1

// Application Gateway
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGatewayPrivateLinkConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/privateLinkConfigurations/privateLinkConfiguration1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontendPort -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/frontendPorts/feport1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontendIPConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/frontendIPConfigurations/feipconfig1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BackendAddressPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/backendAddressPools/beap1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HttpListener -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/httpListeners/listener1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AuthenticationCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/authenticationCertificates/authcert1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RewriteRuleSet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/rewriteRuleSets/rewriteRuleSet1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Probe -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/probes/probe1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SslCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/sslCertificates/sslcert1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BackendHttpSettingsCollection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/backendHttpSettingsCollection/backendHttpSettingsCollection1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RedirectConfigurations -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/redirectConfigurations/redirectConfig1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=TrustedRootCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/trustedRootCertificates/rootCert1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=UrlPathMap -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/urlPathMaps/urlpath1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SslProfile -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/sslProfiles/sslprofile1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=TrustedClientCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/applicationGateways/applicationGateway1/trustedClientCertificates/trustedClientCert1 -rewrite=true

// NAT Gateway
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NatGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/natGateways/gateway1
// NOTE: the Nat Gateway <-> Public IP Association can't be generated at this time

// Private Link
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateDnsZoneConfig -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1/privateDnsZoneConfigs/privateDnsZoneConfig1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateDnsZoneGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1
// ^ these two looks like it should be in Private DNS - alas no, it's actually nested and entirely managed within the Private Endpoint

// Virtual Hubs
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BgpConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/bgpConnections/connection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HubRouteTable -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/routeTable1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HubRouteTableRoute -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/routeTable1/routes/route1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HubVirtualNetworkConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubVirtualNetworkConnections/hubConnection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SecurityPartnerProvider -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/securityPartnerProviders/partnerProvider1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualHub -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualHubIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/ipConfigurations/ipConfiguration1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualWan -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualWans/virtualWan1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PointToSiteVpnGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/p2sVpnGateways/pointToSite1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VpnConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1/vpnConnections/vpnConnection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VpnGatewayNatRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1/natRules/natRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RouteMap -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/vhub1/routeMaps/routeMap1

// Subnet Service Endpoint Policy
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SubnetServiceEndpointStoragePolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/serviceEndpointPolicies/policy1

// Express Route Port
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRoutePort -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRoutePorts/port1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRoutePortAuthorization -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRoutePorts/port1/authorizations/authorization1

// Virtual Network Gateway
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkGatewayIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/ipConfigurations/cfg1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkGatewayNatRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/natRules/rule1

// Express Route Connection
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteCircuit -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/erCircuit1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteCircuitPeering -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/erCircuit1/peerings/peering1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteGateways/ergw1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteCircuitConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/connection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteGateways/ergw1/expressRouteConnections/erConnection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExpressRouteCircuitAuthorization -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/expressRouteCircuit1/authorizations/authorization1

// Network
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkInterfaceIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1/ipConfigurations/config1

// Virtual Machine Scale Set
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualMachineScaleSetPublicIPAddress -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfiguration1/publicIPAddresses/publicIpAddress1

// Custom IP Prefix
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CustomIpPrefix -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/customIPPrefixes/prefix1

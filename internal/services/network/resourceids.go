// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

// Core bits and pieces
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGatewayURLPathMapPathRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/applicationGateways/applicationGateway1/urlPathMaps/urlPathMap1/pathRules/pathRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkDnsServers -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/dnsServers/default -rewrite=true

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

// Private Link
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateDnsZoneConfig -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1/privateDnsZoneConfigs/privateDnsZoneConfig1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateDnsZoneGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1/privateDnsZoneGroups/privateDnsZoneGroup1
// ^ these two looks like it should be in Private DNS - alas no, it's actually nested and entirely managed within the Private Endpoint

// Virtual Hubs
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HubRouteTableRoute -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/routeTable1/routes/route1

// Virtual Network Gateway
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkGatewayIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/ipConfigurations/cfg1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkGatewayPolicyGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/virtualNetworkGatewayPolicyGroups/policyGroup1

// Network
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkInterfaceIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1/ipConfigurations/config1

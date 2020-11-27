package network

// Core bits and pieces
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IpGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/ipGroups/group1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkInterface -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PublicIpAddress -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/publicIPAddresses/publicIpAddress1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Subnet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/subnets/subnet1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetwork -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1

// Firewall Policy
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallPolicyRuleCollectionGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/ruleCollectionGroup1

// Load Balancers
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancer -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerBackendAddressPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerFrontendIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/frontendIPConfigurations/frontendIPConfig1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerInboundNatPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/inboundNatPools/pool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerInboundNatRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/inboundNatRules/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerOutboundRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/outboundRules/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerProbe -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/probes/probe1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancingRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/loadBalancingRules/rule1

// NAT Gateway
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NatGateway -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/natGateways/gateway1
// NOTE: the Nat Gateway <-> Public IP Association can't be generated at this time

// Network Watcher
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ConnectionMonitor -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkWatchers/watcher1/connectionMonitors/connectionMonitor1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NetworkWatcher -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkWatchers/watcher1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PacketCapture -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkWatchers/watcher1/packetCaptures/capture1

// Private Link
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateEndpoint -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateEndpoints/endpoint1

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

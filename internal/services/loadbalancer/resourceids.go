// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

// Load Balancers
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BackendAddressPoolAddress -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/addresses/address1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerBackendAddressPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerFrontendIpConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/frontendIPConfigurations/frontendIPConfig1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerInboundNatPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/inboundNatPools/pool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancerInboundNatRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/inboundNatRules/rule1
